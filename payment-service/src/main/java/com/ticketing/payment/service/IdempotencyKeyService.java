package com.ticketing.payment.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.ticketing.payment.entity.IdempotencyKey;
import com.ticketing.payment.entity.enums.IdempotencyStatus;
import com.ticketing.payment.repository.IdempotencyKeyRepository;
import com.ticketing.payment.service.dto.IdempotencyKeyContext;
import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.time.LocalDateTime;
import java.util.Map;
import java.util.Optional;

/**
 * Handles creation and retrieval of idempotency keys to prevent duplicate
 * payment/refund requests.
 */
@Service
@RequiredArgsConstructor
@Slf4j
public class IdempotencyKeyService {

    private final IdempotencyKeyRepository idempotencyKeyRepository;
    private final ObjectMapper objectMapper;

    /**
     * Ensure an idempotency key exists before processing a request. If key is
     * found and completed, returns stored response. If still processing,
     * instruct caller to wait/abort. Otherwise create new key entry.
     */
    @Transactional
    public Optional<IdempotencyKey> registerOrRetrieve(IdempotencyKeyContext context) {
        if (!StringUtils.hasText(context.getKey())) {
            return Optional.empty();
        }

        Optional<IdempotencyKey> existing = idempotencyKeyRepository.findByIdempotencyKey(context.getKey());
        if (existing.isPresent()) {
            IdempotencyKey key = existing.get();
            if (key.isCompleted()) {
                log.info("Idempotency key {} already completed, returning cached response", key.getIdempotencyKey());
                return existing;
            }
            if (key.isProcessing() && !key.isExpired()) {
                throw new IllegalStateException("Request is still processing for this idempotency key");
            }
            // Expired or failed -> update metadata and reuse
            updateExistingKey(key, context);
            return Optional.of(key);
        }

        IdempotencyKey key = createKey(context);
        idempotencyKeyRepository.save(key);
        return Optional.of(key);
    }

    /**
     * Mark key as completed, storing response payload for future replay.
     */
    @Transactional
    public void markCompleted(IdempotencyKey key, int responseStatus, Map<String, Object> responseBody) {
        key.markAsCompleted(responseStatus, writeJson(responseBody));
        idempotencyKeyRepository.save(key);
    }

    /**
     * Mark key as failed with optional payload.
     */
    @Transactional
    public void markFailed(IdempotencyKey key, int responseStatus, Map<String, Object> responseBody) {
        key.markAsFailed(responseStatus, writeJson(responseBody));
        idempotencyKeyRepository.save(key);
    }

    /**
     * Cleanup API: remove keys older than cutoff.
     */
    @Transactional
    public long purgeKeysOlderThan(LocalDateTime cutoff) {
        long removed = idempotencyKeyRepository.deleteByCreatedAtBefore(cutoff);
        log.info("Purged {} idempotency keys older than {}", removed, cutoff);
        return removed;
    }

    private IdempotencyKey createKey(IdempotencyKeyContext context) {
        return IdempotencyKey.builder()
                .idempotencyKey(context.getKey())
                .requestPath(context.getRequestPath())
                .requestMethod(context.getRequestMethod())
                .userId(context.getUserId())
                .ipAddress(context.getIpAddress())
                .userAgent(context.getUserAgent())
                .requestHeaders(writeJson(context.getRequestHeaders()))
                .requestBody(writeJson(context.getRequestBody()))
                .expiresAt(LocalDateTime.now().plusHours(24))
                .status(IdempotencyStatus.PROCESSING)
                .build();
    }

    private void updateExistingKey(IdempotencyKey key, IdempotencyKeyContext context) {
        key.setRequestPath(context.getRequestPath());
        key.setRequestMethod(context.getRequestMethod());
        key.setUserId(context.getUserId());
        key.setIpAddress(context.getIpAddress());
        key.setUserAgent(context.getUserAgent());
        key.setRequestHeaders(writeJson(context.getRequestHeaders()));
        key.setRequestBody(writeJson(context.getRequestBody()));
        key.setExpiresAt(LocalDateTime.now().plusHours(24));
        key.setStatus(IdempotencyStatus.PROCESSING);
        key.setResponseStatus(null);
        key.setResponseBody(null);
        key.setResponseHeaders(null);
        idempotencyKeyRepository.save(key);
    }

    private String writeJson(Map<String, Object> data) {
        if (data == null) {
            return null;
        }
        try {
            return objectMapper.writeValueAsString(data);
        } catch (JsonProcessingException e) {
            throw new IllegalArgumentException("Failed to serialize idempotency payload", e);
        }
    }
}

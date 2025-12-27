package com.ticketing.payment.service;

import com.ticketing.payment.entity.Payment;
import com.ticketing.payment.entity.Refund;
import com.ticketing.payment.entity.TransactionLog;
import com.ticketing.payment.entity.enums.PaymentGateway;
import com.ticketing.payment.entity.enums.TransactionStatus;
import com.ticketing.payment.entity.enums.TransactionType;
import com.ticketing.payment.repository.TransactionLogRepository;
import jakarta.transaction.Transactional;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.UUID;

/**
 * Service responsible for persisting transaction logs and providing helper
 * operations for monitoring/cleanup.
 */
@Service
@RequiredArgsConstructor
@Slf4j
public class TransactionLogService {

    private final TransactionLogRepository transactionLogRepository;

    @Transactional
    public void logPaymentEvent(TransactionType type, TransactionStatus status, Payment payment, String eventName,
            String externalReference, String errorCode, String errorMessage) {
        TransactionLog txLog = TransactionLog.builder()
                .payment(payment)
                .paymentUuid(payment.getPaymentId())
                .transactionType(type)
                .eventName(eventName)
                .status(status)
                .gatewayProvider(payment.getGatewayProvider())
                .externalReference(externalReference)
                .errorCode(errorCode)
                .errorMessage(errorMessage)
                .build();
        transactionLogRepository.save(txLog);
        log.debug("Payment transaction log recorded: {}", txLog.getLogId());
    }

    @Transactional
    public void logRefundEvent(TransactionType type, TransactionStatus status, Refund refund, String eventName,
            String externalReference, String errorCode, String errorMessage) {
        TransactionLog txLog = TransactionLog.builder()
                .payment(refund.getPayment())
                .paymentUuid(refund.getPaymentUuid())
                .refund(refund)
                .refundUuid(refund.getRefundId())
                .transactionType(type)
                .eventName(eventName)
                .status(status)
                .gatewayProvider(refund.getGatewayProvider())
                .externalReference(externalReference)
                .errorCode(errorCode)
                .errorMessage(errorMessage)
                .build();
        transactionLogRepository.save(txLog);
    }

    /**
     * Cleanup helper for retention policy.
     */
    @Transactional
    public long purgeLogsOlderThan(LocalDateTime cutoff) {
        long removed = transactionLogRepository.deleteByCreatedAtBefore(cutoff);
        log.info("Purged {} transaction logs older than {}", removed, cutoff);
        return removed;
    }
}

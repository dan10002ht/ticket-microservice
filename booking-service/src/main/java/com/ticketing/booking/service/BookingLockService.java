package com.ticketing.booking.service;

import java.util.concurrent.TimeUnit;

import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import lombok.extern.slf4j.Slf4j;

/**
 * Distributed lock service for booking operations.
 * Uses Redisson (Redis-based) for distributed locking.
 *
 * IMPORTANT: Lock lease time must be longer than saga duration to prevent
 * concurrent saga execution on the same event.
 */
@Slf4j
@Component
public class BookingLockService {

    private final RedissonClient redissonClient;
    private final Counter lockAcquiredCounter;
    private final Counter lockFailedCounter;

    /**
     * Wait time to acquire lock (how long to wait if lock is held by another thread)
     * Default: 10 seconds
     */
    @Value("${booking.lock.wait-time-seconds:10}")
    private long lockWaitTimeSeconds;

    /**
     * Lease time (how long the lock is held before auto-release)
     * MUST be longer than the maximum saga execution time (typically 30-60 seconds)
     * Default: 120 seconds (2 minutes) to account for network latency and retries
     */
    @Value("${booking.lock.lease-time-seconds:120}")
    private long lockLeaseTimeSeconds;

    public BookingLockService(RedissonClient redissonClient, MeterRegistry meterRegistry) {
        this.redissonClient = redissonClient;
        this.lockAcquiredCounter = meterRegistry.counter("booking.lock.acquired");
        this.lockFailedCounter = meterRegistry.counter("booking.lock.failed");
    }

    /**
     * Acquire a distributed lock for booking operations.
     *
     * @param key The lock key (typically event ID)
     * @return The acquired lock
     * @throws InterruptedException if interrupted while waiting
     * @throws IllegalStateException if lock cannot be acquired
     */
    public RLock acquireLock(String key) throws InterruptedException {
        RLock lock = redissonClient.getLock("booking:lock:" + key);

        // tryLock(waitTime, leaseTime, unit)
        // - waitTime: how long to wait to acquire the lock
        // - leaseTime: how long to hold the lock before auto-release
        boolean acquired = lock.tryLock(lockWaitTimeSeconds, lockLeaseTimeSeconds, TimeUnit.SECONDS);

        if (!acquired) {
            lockFailedCounter.increment();
            log.warn("Failed to acquire lock for key: {}, waitTime: {}s", key, lockWaitTimeSeconds);
            throw new IllegalStateException("Unable to acquire booking lock for key " + key);
        }

        lockAcquiredCounter.increment();
        log.debug("Lock acquired for key: {}, leaseTime: {}s", key, lockLeaseTimeSeconds);
        return lock;
    }

    /**
     * Release a lock if held by current thread.
     *
     * @param lock The lock to release
     */
    public void releaseLock(RLock lock) {
        if (lock != null && lock.isHeldByCurrentThread()) {
            try {
                lock.unlock();
                log.debug("Lock released");
            } catch (IllegalMonitorStateException e) {
                // Lock was already released (possibly due to lease expiry)
                log.warn("Lock was already released, possibly due to lease expiry");
            }
        }
    }

    /**
     * Extend lock lease time if needed for long-running operations.
     * Call this if saga is taking longer than expected.
     *
     * @param lock The lock to extend
     * @param additionalSeconds Additional seconds to extend
     * @return true if extended successfully
     */
    public boolean extendLock(RLock lock, long additionalSeconds) {
        if (lock != null && lock.isHeldByCurrentThread()) {
            try {
                // Redisson automatically extends lock if we're still holding it
                // But we can manually check remaining TTL
                long remainingTTL = lock.remainTime();
                if (remainingTTL > 0) {
                    log.debug("Lock still valid, remaining TTL: {}ms", remainingTTL);
                    return true;
                }
            } catch (Exception e) {
                log.error("Failed to check/extend lock", e);
            }
        }
        return false;
    }
}


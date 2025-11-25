package com.ticketing.booking.service;

import java.util.concurrent.TimeUnit;

import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import lombok.RequiredArgsConstructor;

@Component
@RequiredArgsConstructor
public class BookingLockService {

    private final RedissonClient redissonClient;

    @Value("${booking.lock.timeout-seconds:5}")
    private long lockTimeoutSeconds;

    public RLock acquireLock(String key) throws InterruptedException {
        RLock lock = redissonClient.getLock("booking:lock:" + key);
        boolean acquired = lock.tryLock(lockTimeoutSeconds, TimeUnit.SECONDS);
        if (!acquired) {
            throw new IllegalStateException("Unable to acquire booking lock for key " + key);
        }
        return lock;
    }

    public void releaseLock(RLock lock) {
        if (lock != null && lock.isHeldByCurrentThread()) {
            lock.unlock();
        }
    }
}


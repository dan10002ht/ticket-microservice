package com.ticketing.payment.adapter;

import com.ticketing.payment.entity.enums.PaymentGateway;

import java.util.EnumMap;
import java.util.Map;

public class PaymentGatewayRegistry {

    private final Map<PaymentGateway, PaymentGatewayAdapter> adapters = new EnumMap<>(PaymentGateway.class);

    public void register(PaymentGateway gateway, PaymentGatewayAdapter adapter) {
        adapters.put(gateway, adapter);
    }

    public PaymentGatewayAdapter resolve(PaymentGateway gateway) {
        PaymentGatewayAdapter adapter = adapters.get(gateway);
        if (adapter == null) {
            throw new IllegalArgumentException("No adapter registered for gateway " + gateway);
        }
        return adapter;
    }
}

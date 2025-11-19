package com.ticketing.payment.adapter;

import com.ticketing.payment.service.dto.PaymentCreateCommand;
import com.ticketing.payment.service.dto.RefundCreateCommand;

/**
 * Base interface for all payment gateway implementations.
 */
import com.ticketing.payment.entity.enums.PaymentGateway;

public interface PaymentGatewayAdapter {

    PaymentGateway getGateway();

    GatewayResponse authorize(PaymentCreateCommand command, GatewayRequestContext context);

    GatewayResponse capture(PaymentCreateCommand command, GatewayRequestContext context);

    GatewayResponse refund(RefundCreateCommand command, GatewayRequestContext context);

    GatewayResponse cancelPayment(String providerReference, String reason, GatewayRequestContext context);
}

package com.example.ordermgmt.activities;

import com.example.ordermgmt.model.OrderInput;
import com.example.ordermgmt.model.OrderItem;
import io.temporal.activity.ActivityInterface;
import io.temporal.activity.ActivityMethod;
import io.temporal.activity.ActivityOptions;
import io.temporal.activity.LocalActivityOptions;
import io.temporal.common.RetryOptions;

import java.time.Duration;
import java.util.List;

@ActivityInterface
public interface OrderActivities {

    ActivityOptions defaultActivityOptions = ActivityOptions.newBuilder()
            .setStartToCloseTimeout(Duration.ofSeconds(5))
            .setRetryOptions(RetryOptions.newBuilder()
                    .setInitialInterval(Duration.ofSeconds(1))
                    .setBackoffCoefficient(2)
                    .setMaximumInterval(Duration.ofSeconds(30))
                    .build())
            .build();

    LocalActivityOptions defaultLocalActivityOptions =
            LocalActivityOptions.newBuilder()
                    .setStartToCloseTimeout(Duration.ofSeconds(5))
                    .build();

    @ActivityMethod
    List<OrderItem> getItems();

    @ActivityMethod
    String checkFraud(OrderInput input);

    @ActivityMethod
    String prepareShipment(OrderInput input);

    @ActivityMethod
    String chargeCustomer(OrderInput input, String type);

    @ActivityMethod
    void shipOrder(OrderInput input, OrderItem item);

    @ActivityMethod
    String undoPrepareShipment(OrderInput input);

    @ActivityMethod
    String undoChargeCustomer(OrderInput input);

    @ActivityMethod
    String getShippingTaskQueue();

    @ActivityMethod
    public String getShippingServiceEndpoint();
}

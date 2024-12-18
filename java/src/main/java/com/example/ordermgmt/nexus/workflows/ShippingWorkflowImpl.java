package com.example.ordermgmt.nexus.workflows;

import com.example.ordermgmt.activities.OrderActivities;
import com.example.ordermgmt.model.OrderInput;
import com.example.ordermgmt.model.OrderItem;
import com.example.ordermgmt.model.OrderOutput;
import com.example.ordermgmt.model.ShippingInput;
import io.temporal.workflow.Workflow;
import org.slf4j.Logger;

public class ShippingWorkflowImpl implements ShippingWorkflow {
    private static final Logger log = Workflow.getLogger(com.example.ordermgmt.workflows.ShippingChildWorkflowImpl.class);

    private final OrderActivities activities = Workflow.newActivityStub(OrderActivities.class,
            OrderActivities.defaultActivityOptions);


    @Override
    public OrderOutput execute(ShippingInput input) {


        log.info("Shipping workflow started, orderId = {}", input.getOrderInput().getOrderId());
        // Ship order
        activities.shipOrder(input.getOrderInput(), input.getOrderItem());

        return new OrderOutput("OrderTrackingID", "");
    }
}

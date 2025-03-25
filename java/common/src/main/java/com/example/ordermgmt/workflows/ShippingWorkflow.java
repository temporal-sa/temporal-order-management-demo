package com.example.ordermgmt.workflows;

import com.example.ordermgmt.activities.OrderActivities;
import com.example.ordermgmt.model.OrderOutput;
import com.example.ordermgmt.model.ShippingInput;
import io.temporal.spring.boot.WorkflowImpl;
import io.temporal.workflow.Workflow;
import io.temporal.workflow.WorkflowInterface;
import io.temporal.workflow.WorkflowMethod;
import org.slf4j.Logger;

@WorkflowInterface
public interface ShippingWorkflow {
    @WorkflowMethod(name = "ShippingWorkflow")
    public OrderOutput execute(ShippingInput input);


    @WorkflowImpl
    class ShippingWorkflowImpl implements ShippingWorkflow {
        private static final Logger log = Workflow.getLogger(ShippingWorkflowImpl.class);

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
}


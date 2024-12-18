package com.example.ordermgmt.nexus.workflows;

import com.example.ordermgmt.model.OrderOutput;
import com.example.ordermgmt.model.ShippingInput;
import io.temporal.workflow.WorkflowInterface;
import io.temporal.workflow.WorkflowMethod;

@WorkflowInterface
public interface ShippingWorkflow {
    @WorkflowMethod
    public OrderOutput execute(ShippingInput input);
}


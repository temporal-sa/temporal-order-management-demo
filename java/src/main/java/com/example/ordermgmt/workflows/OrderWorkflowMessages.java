package com.example.ordermgmt.workflows;

import com.example.ordermgmt.model.UpdateOrderInput;
import io.temporal.workflow.QueryMethod;
import io.temporal.workflow.SignalMethod;
import io.temporal.workflow.UpdateMethod;
import io.temporal.workflow.UpdateValidatorMethod;

public interface OrderWorkflowMessages {
    @QueryMethod(name = "getProgress")
    int queryProgress();

    @SignalMethod(name = "UpdateOrder")
    void updateOrderSignal(UpdateOrderInput updateInput);

    @UpdateMethod(name = "UpdateOrder")
    String updateOrderUpdate(UpdateOrderInput updateInput);

    @UpdateValidatorMethod(updateName = "UpdateOrder")
    void updateOrderValidator(UpdateOrderInput updateInput);
}

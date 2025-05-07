package com.example.ordermgmt.nexus.handler;

import com.example.ordermgmt.model.ShippingInput;
import io.nexusrpc.Operation;
import io.nexusrpc.Service;

@Service
public interface ShippingService {
    @Operation
    String execute(ShippingInput input);
}

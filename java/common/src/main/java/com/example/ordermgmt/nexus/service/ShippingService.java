package com.example.ordermgmt.nexus.service;

import com.example.ordermgmt.model.OrderInput;
import com.example.ordermgmt.model.OrderItem;
import com.example.ordermgmt.model.OrderOutput;
import com.example.ordermgmt.model.ShippingInput;
import io.nexusrpc.Operation;
import io.nexusrpc.Service;

@Service
public interface ShippingService {

    @Operation
    OrderOutput execute(ShippingInput input);

}

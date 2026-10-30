package com.example.ordermgmt.model;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ShippingInput {
    OrderInput orderInput;
    OrderItem orderItem;
    int index; // item n ships in n*2 seconds; 0 means a random delay

    public ShippingInput(OrderInput orderInput, OrderItem orderItem) {
        this(orderInput, orderItem, 0);
    }
}

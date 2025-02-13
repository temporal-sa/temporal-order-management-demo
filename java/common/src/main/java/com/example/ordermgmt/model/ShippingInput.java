package com.example.ordermgmt.model;

import lombok.*;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ShippingInput {
    OrderInput orderInput;
    OrderItem orderItem;
}

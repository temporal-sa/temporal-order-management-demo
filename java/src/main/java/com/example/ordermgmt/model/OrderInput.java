package com.example.ordermgmt.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class OrderInput {
    @JsonProperty("OrderId")
    private String orderId;

    @JsonProperty("Address")
    private String address;
}

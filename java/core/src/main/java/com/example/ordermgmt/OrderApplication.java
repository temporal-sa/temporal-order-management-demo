package com.example.ordermgmt;

import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@Slf4j
@SpringBootApplication
public class OrderApplication {
    public static final String defaultShippingNexusEndpoint =
            System.getenv("TEMPORAL_NEXUS_SHIPPING_ENDPOINT") != null ?
                    System.getenv("TEMPORAL_NEXUS_SHIPPING_ENDPOINT") :
                    "shipping-endpoint";

    public static void main(String[] args) {
        SpringApplication.run(OrderApplication.class, args);
    }

}

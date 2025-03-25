package com.example.ordermgmt.customize;

import io.temporal.workflow.NexusServiceOptions;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import io.temporal.client.WorkflowClientOptions;
import io.temporal.common.converter.CodecDataConverter;
import io.temporal.common.converter.DefaultDataConverter;
import io.temporal.serviceclient.WorkflowServiceStubsOptions;
import io.temporal.spring.boot.TemporalOptionsCustomizer;
import io.temporal.spring.boot.WorkerOptionsCustomizer;
import io.temporal.worker.WorkerFactoryOptions;
import io.temporal.worker.WorkerOptions;
import io.temporal.worker.WorkflowImplementationOptions;
import io.temporal.worker.tuning.ResourceBasedControllerOptions;
import io.temporal.worker.tuning.ResourceBasedTuner;
import jakarta.annotation.Nonnull;

import java.lang.management.ManagementFactory;
import java.util.Collections;


@Configuration
public class TemporalOptionsConfig {

    @Value("${nexus.worker.endpoint}")
    private String shippingEndpoint;

    @Bean
    public WorkerOptionsCustomizer customWorkerOptions() {
        return new WorkerOptionsCustomizer() {
            @Nonnull
            @Override
            public WorkerOptions.Builder customize(
                    @Nonnull WorkerOptions.Builder optionsBuilder,
                    @Nonnull String workerName,
                    @Nonnull String taskQueue) {

                // Adding the tuning option for autotuning the worker to use 75% of the memory and CPU available.
                optionsBuilder.setIdentity(taskQueue+"-"+ ManagementFactory.getRuntimeMXBean().getName())
                              .setWorkerTuner(
                        ResourceBasedTuner.newBuilder()
                                .setControllerOptions(ResourceBasedControllerOptions.newBuilder(0.75, 0.75)
                                        .build())
                              .build());
                return optionsBuilder;
            }
        };
    }

    // WorkflowServiceStubsOptions customization
    @Bean
    public TemporalOptionsCustomizer<WorkflowServiceStubsOptions.Builder>
    customServiceStubsOptions() {
        return new TemporalOptionsCustomizer<WorkflowServiceStubsOptions.Builder>() {
            @Nonnull
            @Override
            public WorkflowServiceStubsOptions.Builder customize(
                    @Nonnull WorkflowServiceStubsOptions.Builder optionsBuilder) {
                // set options on optionsBuilder as needed
                // ...
                return optionsBuilder;
            }
        };
    }

    // WorkflowClientOption customization
    // Only option we are adding just now is for the codec converter.
    @Bean
    public TemporalOptionsCustomizer<WorkflowClientOptions.Builder> customClientOptions() {
        return new TemporalOptionsCustomizer<WorkflowClientOptions.Builder>() {
            @Nonnull
            @Override
            public WorkflowClientOptions.Builder customize(
                    @Nonnull WorkflowClientOptions.Builder optionsBuilder) {


                return optionsBuilder;
            }
        };
    }

    // WorkerFactoryOptions customization
    @Bean
    public TemporalOptionsCustomizer<WorkerFactoryOptions.Builder> customWorkerFactoryOptions() {
        return new TemporalOptionsCustomizer<WorkerFactoryOptions.Builder>() {
            @Nonnull
            @Override
            public WorkerFactoryOptions.Builder customize(
                    @Nonnull WorkerFactoryOptions.Builder optionsBuilder) {
                // set options on optionsBuilder as needed
                // ...
                return optionsBuilder;
            }
        };
    }

    // WorkflowImplementationOptions customization
    @Bean
    public TemporalOptionsCustomizer<WorkflowImplementationOptions.Builder>
    customWorkflowImplementationOptions() {
        return new TemporalOptionsCustomizer<>() {
            @Nonnull
            @Override
            public WorkflowImplementationOptions.Builder customize(
                    @Nonnull WorkflowImplementationOptions.Builder optionsBuilder) {
                // set options on optionsBuilder such as per-activity options

                optionsBuilder.setNexusServiceOptions(
                        Collections.singletonMap(
                                "ShippingService",
                                NexusServiceOptions.newBuilder().setEndpoint(shippingEndpoint).build()
                        )
                ).build();
                return optionsBuilder;
            }
        };
    }
}


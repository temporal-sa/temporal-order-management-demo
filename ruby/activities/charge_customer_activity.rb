require 'temporalio/activity'

module Activities
  class ChargeCustomerActivity < Temporalio::Activity::Definition
    ERROR_CHARGE_API_UNAVAILABLE = 'OrderWorkflowAPIFailure'
    ERROR_INVALID_CREDIT_CARD = 'OrderWorkflowNonRecoverableFailure'

    def execute(_input, type)
      attempt = Temporalio::Activity::Context.current.info.attempt
      error = simulate_external_operation_charge(1000, type, attempt)
      case error
      when ERROR_CHARGE_API_UNAVAILABLE
        raise StandardError.new('Charge Customer activity failed, API unavailable')
      when ERROR_INVALID_CREDIT_CARD
        raise Temporalio::Error::ApplicationError.new('Charge Customer activity failed, card is invalid', non_retryable: true)
      end
      'SUCCESS'
    end

    private

    def simulate_external_operation(ms)
      sleep(ms / 1000.0)
    end

    def simulate_external_operation_charge(ms, type, attempt)
      simulate_external_operation(ms / attempt)
      attempt < 5 ? type : 'NoError'
    end
  end

  def logger
    @logger ||= Temporalio::Workflow.logger
  end
end

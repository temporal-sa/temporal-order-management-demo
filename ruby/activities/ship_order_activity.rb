require 'temporalio/activity'

module Activities
  class ShipOrderActivity < Temporalio::Activity::Definition
    def execute(input, item)
      delay_ms = rand(1000..4000)
      sleep(delay_ms / 1000.0)
      'SUCCESS'
    end
  end
end

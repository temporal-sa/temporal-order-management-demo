require 'temporalio/activity'

module Activities
  class PrepareShipmentActivity < Temporalio::Activity::Definition
    def execute(input)
      sleep(1)
      'SUCCESS'
    end
  end
end 
require 'temporalio/activity'

module Activities
  class CheckFraudActivity < Temporalio::Activity::Definition
    def execute(input)
      sleep(1)
      'SUCCESS'
    end
  end
end 
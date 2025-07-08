require 'temporalio/activity'
require_relative '../models/order_item'

module Activities
  class GetItemsActivity < Temporalio::Activity::Definition
    def execute
      [
        Models::OrderItem.new(654300, "Table Top", 1),
        Models::OrderItem.new(654321, "Table Legs", 2),
        Models::OrderItem.new(654322, "Keypad", 1)
      ]
    end
  end
end 
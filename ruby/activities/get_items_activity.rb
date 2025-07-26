require 'temporalio/activity'
require_relative '../models/order_item'

module Activities
  class GetItemsActivity < Temporalio::Activity::Definition
    def execute
      [
        Models::OrderItem.new(654_300, 'Table Top', 1),
        Models::OrderItem.new(654_321, 'Table Legs', 2),
        Models::OrderItem.new(654_322, 'Keypad', 1)
      ]
    end
  end
end

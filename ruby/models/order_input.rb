require 'json/add/struct'
require_relative 'serialization'

module Models
  OrderInput = Struct.new(:order_id, :address) do
    include Serialization
  end
end

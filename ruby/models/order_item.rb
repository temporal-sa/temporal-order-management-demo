require 'json/add/struct'
require_relative 'serialization'

module Models
  OrderItem = Struct.new(:id, :description, :quantity) do
    include Serialization
  end
end

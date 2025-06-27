require 'json/add/struct'
require_relative 'serialization'

module Models
  ActivityResponse = Struct.new(:order_id, :status) do
    include Serialization
  end
end 
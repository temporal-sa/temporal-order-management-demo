require 'json'

module Models
  module Serialization
    def self.included(base)
      base.class_eval do
        def to_h
          hash = {}
          members.each do |member|
            value = send(member)
            hash[member] = if value.respond_to?(:to_h)
                             value.to_h
                           else
                             value
                           end
          end
          hash
        end
      end
    end
  end
end 
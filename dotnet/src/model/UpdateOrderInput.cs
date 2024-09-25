using System.Text.Json.Serialization;

namespace DotNetOrderManagement.model;

public class UpdateOrderInput
{
    [JsonPropertyName("Address")]
    public string Address { get; set; } 

    public UpdateOrderInput()
    {
        Address = string.Empty;
    }

    public UpdateOrderInput(string address)
    {
        Address = address;
    }
    
}
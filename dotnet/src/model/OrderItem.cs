namespace DotNetOrderManagement.model;

public class OrderItem
{
    public int Id { get; set; }
    public String Description { get; set; }
    public int Quantity { get; set; }

    public OrderItem()
    {
        Id = 0;
        Description = string.Empty;
        Quantity = 0;
    }

    public OrderItem(int id, String description, int quantity)
    {
        this.Id = id;
        this.Description = description;
        this.Quantity = quantity;
    }
}
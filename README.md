Design Patterns Used
    Builder Pattern - Document, Table, Cell, Config builders
    Factory Pattern - DocumentFactory, FormBuilder
    Functional Options - Client configuration
    Strategy Pattern - RetryPolicy interface

Usage as External Module

```code

import pdf "github.com/chinmay/gopdfsuit-client"

// Read from JSON
client := pdf.NewClient("http://localhost:8080")
doc, _ := client.ReadFromFile(ctx, "form.json")
client.SendAndSave(ctx, doc, "output.pdf")

// Build with builder pattern
doc := pdf.NewDocumentBuilder().
    WithConfig(config).
    WithTitle("font1:16:100:center:0:0:0:1", "My Form").
    AddTable(table).
    Build()

// Quick form with factory
doc := pdf.NewFormBuilder().
    WithTitle("Contact Form").
    AddTextField("Name:", "name", "John").
    Build()

```
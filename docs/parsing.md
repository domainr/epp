# EPP Response Parsing Architecture

This document explains the architecture used to parse EPP XML responses, particularly within `check.go`. Understanding this pattern is crucial for extending the library to support new EPP extensions.

## The Challenge: Diverse and Inconsistent XML

The EPP standard allows for various extensions, which registries use to provide non-standard data, such as domain premium pricing. Over the years, many different versions of these extensions have been created (e.g., `fee-0.5`, `fee-1.0`, `charge-1.0`, `namestore`, etc.).

These extensions result in XML responses that are structurally different and often inconsistent. Using Go's standard `xml.Unmarshal` with a single, static struct would be nearly impossible to maintain. The struct would become enormous, and handling the subtle differences between extension versions would require complex logic.

## The Solution: An Event-Driven Scanner

To solve this, our library uses a flexible, event-driven XML scanner provided by the `github.com/nbio/xx` package.

Instead of trying to unmarshal the entire XML document at once, the scanner walks the XML tree and triggers **handler functions** whenever it enters or reads data from an element that matches a specific path.

### Core Concepts

1.  **`scanResponse`**: This is the global instance of the `xx.Scanner`. It holds all the rules for parsing.

2.  **Path-Based Matching**: Handlers are registered for specific XML paths. The path is a space-separated string representing the hierarchy of XML tags. For example, the path `epp > response > resData > domain:chkData > cd > name` matches the `<name>` element in the following structure:

    ```xml
    <epp>
      <response>
        <resData>
          <domain:chkData>
            <cd>
              <name>example.com</name>
            </cd>
          </domain:chkData>
        </resData>
      </response>
    </epp>
    ```
    > **Note:** The `xx` scanner ignores XML namespace prefixes (like `domain:`), so they are not included in the path strings in the code.


### Example Walkthrough: Parsing a `fee-1.0` Extension

Let's trace how the parser handles the `fee-1.0` extension from our test case.

**Sample XML Snippet:**
```xml
<extension>
    <fee:chkData xmlns:fee="urn:ietf:params:xml:ns:epp:fee-1.0">
        <fee:cd avail="1">
            <fee:name premium="true">zero.work</fee:name>
            <fee:class>premium</fee:class>
            <fee:command name="create">
                <fee:fee description="Registration Fee">500.000</fee:fee>
            </fee:command>
        </fee:cd>
    </fee:chkData>
</extension>
```

**Corresponding Handler Registration in `check.go`'s `init()`:**
```go
// (Path construction combines to form the full path)
path = "epp > response > extension > " + ExtFee10 + " chkData"

// 1. A new DomainCharge is created when <cd> is entered.
scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
    dcr := &c.Value.(*Response).DomainCheckResponse
    dcr.Charges = append(dcr.Charges, DomainCharge{})
    return nil
})

// 2. The domain name is extracted from <name>.
scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
    charges := c.Value.(*Response).DomainCheckResponse.Charges
    charge := &charges[len(charges)-1]
    charge.Domain = string(c.CharData)
    return nil
})

// 3. The category is extracted from <class>.
scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
    charges := c.Value.(*Response).DomainCheckResponse.Charges
    charge := &charges[len(charges)-1]
    if string(c.CharData) != "standard" {
        charge.Category = "premium"
    }
    return nil
})

// 4. The category name is extracted from the <fee> element's attribute.
scanResponse.MustHandleCharData(path+">cd>command>fee", func(c *xx.Context) error {
    charges := c.Value.(*Response).DomainCheckResponse.Charges
    charge := &charges[len(charges)-1]
    charge.CategoryName = c.Attr("", "description")
    return nil
})
```

This modular approach allows us to define parsing logic for each extension independently, making the system highly extensible and maintainable.

## How to Add Support for a New Extension

1.  **Obtain a sample XML response** from the registry that uses the new extension.
2.  **Identify the unique XML path** to the data you need to extract.
3.  **Open `check.go`** and navigate to the `init()` function.
4.  **Add a new block of code** to register handlers for your new extension's path. You will typically use `MustHandleStartElement` to initialize a struct and `MustHandleCharData` to populate its fields.
5.  **Add a new test case** to `response_test.go` with your sample XML to validate that your handlers work correctly.

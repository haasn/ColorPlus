// Example program that uses colorplus, this will take an image example.png and convert it to Adobe RGB for proper viewing on Adobe RGB monitors

package main

import (
    "colorplus"
    "image"
    "image/png"
    "os"
)

func main() {
    // PNG images are encoded as sRGB, so we need to first decode the sRGB, then encode to the Adobe RGB's color space (in linear light), and then apply an sRGB curve again
    chain := colorplus.Chain(colorplus.SpacesRGB.GetDecoder(), colorplus.XYZSpace(colorplus.SpaceAdobeRGB98).GetEncoder(), colorplus.SRGBCurve.GetEncoder())

    // Open the example image
    io, err := os.Open("example.png")

    if err != nil {
        panic(err)
    }

    defer io.Close()

    // Decode this to an NRGBA
    n, err := png.Decode(io)
    image := (*image.RGBA)(n.(*image.NRGBA)) // we need it as RGBA, not NRGBA

    if err != nil {
        panic(err)
    }

    // Apply our chain to the image (argument must be *image.RGBA)
    colorplus.ApplyToImage(image, chain)

    // Save the result
    o, err := os.Create("output.png")

    if err != nil {
        panic(err)
    }

    defer o.Close()

    png.Encode(o, image)
}

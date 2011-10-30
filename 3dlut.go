package colorplus

import "math"

// 3D Lookup table, designed for interop with the 3DL2 and 3DLUT specifications
// TODO: Redesign this from scratch to conform better to the Go style
type LUT3D struct {
    Signature  [4]byte          // file signature; must be: “3DL2” as encoded using the ASCII standard (0x33444C32)
    FileVersion int32           // file format version number (currently 2)
    ProgramName [32]byte        // name of the program that created the file
    ProgramVersion int64        // version number of the program that created the file
    InputBitDepth [3]int32      // input bit depth per component (Y,Cb,Cr or B,G,R)
    InputColorEncoding int32    // input color encoding standard (0 = BGR, 1 = YCbCr)
    InputValueRange int32       // value range for the input (0 = Full range (0-255), 1 = Limited (16-235)) [up to v1]
    OutputBitDepth int32        // output bit depth for all components (valid values are 8, 16, 32 and 64)
    OutputColorEncoding int32   // output color encoding standard (0 = BGR, 1 = YCbCr, 2 = XYZ)
    OutputValueRange int32      // value range for the output
    ParametersFileOffset int32  // number of bytes between the beginning of the file and array 'parametersData'
    ParametersSize int32        // size in bytes of the array 'parametersData'
    LutFileOffset int32         // number of bytes between the beginning of the file and array lutData
    LutCompressionMethod int32  // type of compression used if any (0 = none, 1 = LZO, ...)
    LutCompressedSize int32     // size in bytes of the array 'lutData' inside the file, whether compressed or not
    LutUncompressedSize int32   // size in bytes of the array 'lutData' in memory
    InputColorSpace PRIMARIES   // Input color space - this section is strictly optional,
                                // a value of 0 for all means “disabled/unknown”
    OutputColorSpace PRIMARIES  // Output color space - the same rules apply as for input color space.
                                // In addition, if the output is XYZ, this section is to be set to 0 for all values

    // The array 'parametersData' starts 'parametersFileOffset' bytes after the beginning of the file.
    // The array 'lutDataxx' starts 'lutFileOffset' bytes after the beginning of the file.
    // When creating a 3DLUT2 file, 'lutDataxx' should be positioned on a 16384 byte boundary.

    ParametersData []byte      // Byte array with size 'parametersSize' that should in theory contain an exact copy of the
                                // input file with the commands and settings used for creating the 3DLUT2 file. In practice,
                                // this is however more often used for program-specific meta-data, eg. madVR

    LutData LUTDATA             // Array with size lutSizeUncompressed that contains the 3D LUTs output values.
                                // The type used depends on the outputBitDepth field:
                                //    - unsigned byte,  if outputBitDepth =  8
                                //    - unsigned short, if outputBitDepth = 16
                                //    - float,          if outputBitDepth = 32
                                //    - double,         if outputBitDepth = 64

    //   The value ranges are assumed to be as follows:
    //     Full range (integers):    0 to (2∧depth)−1, for example 0–255 and 0–65535
    //     Full range (floats):      0 to 1
    //     Limited range (integers): 16–235, left/right shifted to the correct bit depth, eg. 4096–60160)
    //     Limited range (floats):   (16÷255)≈0.06275 to (235÷255)≈0.92157.

    //   For XYZ output:
    //     XYZ must be appropriately normalized so that the luminosity of white = the range limit.
    //     eg. for limited range 8-bit output, the luminosity of white would be 235. If the value range is full range,
    //     this would mean that the values are essentially capped to 255. As such, using limited range for XYZ is
    //     preferable for integer output. It is strongly recommended that one uses full range floats for XYZ however,
    //     where the luminosity of white would be normalized to 1.0

    //   The offset inside the array is calculated as:
    //     offset = (A + B << (inputBitDepth[0]) + C << (inputBitDepth[1]+inputBitDepth[0])) × 3

    //   The output order inside the array is:
    //     A = lutDataxx[offset]; B = lutDataxx[offset+1]; C = lutDataxx[offset+2]

    //   A, B and C refer to the respective channels of the used encoding (eg. B, G, R or Y', Cb, Cr)

    //   The lutUncompressedSize of the array is calculated as:
    //     lutDim = 3 × 2∧inputBitDepth[0] × 2∧inputBitDepth[1] × 2∧inputBitDepth[2] × outputBitDepth÷8
}

// Primaries information. This is redundant with the existing ColorSpace structure though, consider removing this in future revisions
type PRIMARIES struct {
    PrimaryRedx, PrimaryRedy,     // The x and y refer here to the coordinates on the CIE xy chromaticity diagram,
    PrimaryGreenx, PrimaryGreeny, // Y is assumed to be 1.0 for these primaries
    PrimaryBluex, PrimaryBluey,
    PrimaryWhitex, PrimaryWhitey float64
}

const (
    RangeFull = 0
    RangeLimited = 1

    EncodingBGR = 0
    EncodingYCbCr = 1
    EncodingXYZ = 2
)

// The actual implementation can be either []byte, []uint16, []float32 or []float64
type LUTDATA interface {
    Get(index int) float64
    Set(index int, value float64)
}

type lutByte []byte
type lutShort []uint16
type lutFloat []float32
type lutDouble []float64

// Byte implementation
func (b lutByte) Get(index int) float64 {
    return float64(b[index])
}

func (b lutByte) Set(index int, value float64) {
    b[index] = byte(math.Floor(value))
}

// Short implementation
func (s lutShort) Get(index int) float64 {
    return float64(s[index])
}

func (s lutShort) Set(index int, value float64) {
    s[index] = uint16(math.Floor(value))
}

// Float implementation
func (f lutFloat) Get(index int) float64 {
    return float64(f[index])
}

func (f lutFloat) Set(index int, value float64) {
    f[index] = float32(value)
}

// Double implementation
func (d lutDouble) Get(index int) float64 {
    return d[index]
}

func (d lutDouble) Set(index int, value float64) {
    d[index] = value
}

// 3DLUT creation
func Make3DLUT(InputDepth []int32, OutputDepth int32, InputRange, OutputRange, InputEncoding, OutputEncoding int32, InputSpace, OutputSpace Space) *LUT3D {
    lut := LUT3D{}

    lut.Signature = [4]byte{'3', 'D', 'L', '2'}
    lut.FileVersion = 2
    lut.ProgramName = [32]byte{'C', 'o', 'l', 'o', 'r', 'P', 'l', 'u', 's'}
    lut.ProgramVersion = 100

    lut.InputBitDepth = [3]int32{ InputDepth[0], InputDepth[1], InputDepth[2] }
    lut.InputColorEncoding = InputEncoding
    lut.InputColorSpace = primariesFromSpace(InputSpace)

    lut.OutputBitDepth = OutputDepth
    lut.OutputColorEncoding = OutputEncoding
    lut.OutputColorSpace = primariesFromSpace(OutputSpace)

    lut.ParametersData = []byte("")
    lut.ParametersFileOffset = 256

    lut.LutUncompressedSize = 3 * (1 << uint(InputDepth[0])) * (1 << uint(InputDepth[1])) * (1 << uint(InputDepth[2])) * (OutputDepth / 8)
    lut.LutCompressionMethod = 0 // LZO currently unsupported
    lut.LutCompressedSize = lut.LutUncompressedSize
    lut.LutFileOffset = 16384

    lut.InputValueRange = InputRange
    lut.OutputValueRange = OutputRange

    lut.LutData = makeLutData(OutputDepth, lut.LutUncompressedSize)

    return &lut
}

// Assignment logic
func (lut *LUT3D) Assign(filter FilterTripleProvider, pipeline bool) {
    var f FilterTriple

    if pipeline {
        f = Chain(Multiplex(
                Pulldown{uint(lut.InputBitDepth[0]), lut.InputValueRange == RangeFull},
                Pulldown{uint(lut.InputBitDepth[1]), lut.InputValueRange == RangeFull},
                Pulldown{uint(lut.InputBitDepth[2]), lut.InputValueRange == RangeFull}),
                filter,
                Pullup{uint(lut.OutputBitDepth), lut.OutputValueRange == RangeFull}).GetTriple()
    } else {
        f = filter.GetTriple()
    }

    // TODO: Actually assign
    lut.SetOutputRaw(0, f(RGB{255, 0, 0}))
}

// Functions for internal logic
func (lut *LUT3D) SetOutputRaw(pos int, output Triple) {
    switch lut.OutputColorEncoding {
    case 0: // BGR
        rgb := output.(RGB)
        lut.LutData.Set(pos, rgb.B)
        lut.LutData.Set(pos + 1, rgb.G)
        lut.LutData.Set(pos + 2, rgb.R)

    case 1: // YCbCr
        panic("[SetOutputRaw] YCbCr currently unsupported!")

    case 2: // XYZ
        var xyz XYZ
        switch x := output.(type) {
            case XYZ: xyz = x
            case Yxy: xyz = x.ToXYZ()
            default: panic("[SetOutputRaw] Unsupported color type for XYZ output")
        }

        lut.LutData.Set(pos, xyz.X)
        lut.LutData.Set(pos + 1, xyz.Y)
        lut.LutData.Set(pos + 2, xyz.Z)

    default: panic("[SetOutputRaw] Unsupported output encoding")
    }
}

func (lut *LUT3D) GetOutputRaw(pos int) Triple {
    a, b, c := lut.LutData.Get(pos), lut.LutData.Get(pos + 1), lut.LutData.Get(pos + 2)

    switch lut.OutputColorEncoding {
        case 0: return RGB{c, b, a} // BGR
        case 1: panic("[GetOutputRaw] YCbCr currently unsupported!") // YCbCr
        case 2: return XYZ{a, b, c} // XYZ
    }

    return nil
}

func (lut *LUT3D) Offset(A, B, C int) int {
    return 3 * ((C << uint(lut.InputBitDepth[1] + lut.InputBitDepth[0])) + (B << uint(lut.InputBitDepth[0])) + A)
}

// Create a new LUTDATA struct based on bit depth
func makeLutData(outdepth, usize int32) LUTDATA {
    switch outdepth {
        case 8: return make(lutByte, usize)
        case 16: return make(lutShort, usize / 2)
        case 32: return make(lutFloat, usize / 4)
        case 64: return make(lutDouble, usize / 8)
        default: panic("[makeLutData] Unsupported output bit depth!")
    }

    return nil
}

// Primaries handling
func primariesFromSpace(space Space) PRIMARIES {
    r, g, b, w := space.Red.ToYxy(), space.Green.ToYxy(), space.Blue.ToYxy(), space.White.ToYxy()

    return PRIMARIES{r.x, r.y, g.x, g.y, b.x, b.y, w.x, w.y}
}

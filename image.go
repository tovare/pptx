package pptx

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"strconv"

	"github.com/beevik/etree"
)

// What should be the image's dimensions?
// When inserting directly, it seems to use 96 dpi, or maybe depending on the presentation setup.
// However optimal resolution should depend on the screen resolution, which is normally unknown.

// An Image can be added to a Slide.
type Image struct {
	X, Y  Dimension
	Image image.Image
}

// addImageFile adds the png file to the presentation.
func (f *File) addImageFile(im Image, imageNum, slideNum int) error {
	imagePath := fmt.Sprintf("ppt/media/slide%dimage%d.png", slideNum, imageNum)
	var b bytes.Buffer
	if err := png.Encode(&b, im.Image); err != nil {
		return fmt.Errorf("slide %d image %d: %s", slideNum, imageNum, err)
	}
	f.m[imagePath] = &b
	return nil
}

// addImageRef adds the image reference to the the slide's xml tree.
// The image reference is appended to the slide at the path:
// <p:sld...><p:cSld><p:spTree>
// addImageRef adds the image reference to the slide's xml tree.
func (s *Slide) addImageRef(im Image, imageNum int) error {
	xml, err := im.build(imageNum)
	if err != nil {
		return err
	}
	root := s.xml.Root()
	if root == nil {
		return fmt.Errorf("Cannot find root element")
	}
	var spTree *etree.Element
	if spTree = root.FindElement("p:cSld/p:spTree"); spTree == nil {
		return fmt.Errorf("Cannot find spTree")
	}
	imRoot := xml.Root()
	if imRoot == nil {
		return fmt.Errorf("Cannot find imRoot element")
	}
	spTree.Child = append(spTree.Child, imRoot)
	return nil
}

// build create the xml tree of the image reference.
func (im *Image) build(imNum int) (*etree.Document, error) {
	cxDim := Dimension(im.Image.Bounds().Dx()) * Inch / 96 // We use 96 dpi to set the images extent.
	cyDim := Dimension(im.Image.Bounds().Dy()) * Inch / 96
	x := strconv.FormatUint(uint64(im.X), 10)
	y := strconv.FormatUint(uint64(im.Y), 10)
	cx := strconv.FormatUint(uint64(cxDim), 10)
	cy := strconv.FormatUint(uint64(cyDim), 10)
	template := `<p:pic>
<p:nvPicPr>
<p:cNvPr id="` + strconv.Itoa(1026+imNum) + `" name="Picture ` + strconv.Itoa(imNum+1) + `"/>
<p:cNvPicPr/>
<p:nvPr/>
</p:nvPicPr>
<p:blipFill>
<a:blip r:embed="rId` + strconv.Itoa(imNum+2) + `">
<a:extLst>
<a:ext uri="{28A0092B-C50C-407E-A947-70E740481C1C}">
<a14:useLocalDpi xmlns:a14="http://schemas.microsoft.com/office/drawing/2010/main" val="0"/>
</a:ext>
</a:extLst>
</a:blip>
<a:srcRect/>
<a:stretch>
<a:fillRect/>
</a:stretch>
</p:blipFill>
<p:spPr bwMode="auto">
<a:xfrm>
<a:off x="` + x + `" y="` + y + `"/>
<a:ext cx="` + cx + `" cy="` + cy + `"/>
</a:xfrm>
<a:prstGeom prst="rect">
<a:avLst/>
</a:prstGeom>
<a:noFill/>
<a:extLst>
<a:ext uri="{909E8E84-426E-40DD-AFC4-6F175D3DCCD1}">
<a14:hiddenFill xmlns:a14="http://schemas.microsoft.com/office/drawing/2010/main">
<a:solidFill>
<a:srgbClr val="FFFFFF"/>
</a:solidFill>
</a14:hiddenFill>
</a:ext>
</a:extLst>
</p:spPr>
</p:pic>`
	doc := etree.NewDocument()
	err := doc.ReadFromString(template)
	return doc, err
}

/* This original image was 192x107
   <p:pic>
     <p:nvPicPr>
       <p:cNvPr id="1026" name="Picture 2"/>
       <p:cNvPicPr/>
       <p:nvPr/>
     </p:nvPicPr>
     <p:blipFill>
       <a:blip r:embed="rId2">
         <a:extLst>
           <a:ext uri="{28A0092B-C50C-407E-A947-70E740481C1C}">
             <a14:useLocalDpi xmlns:a14="http://schemas.microsoft.com/office/drawing/2010/main" val="0"/>
           </a:ext>
         </a:extLst>
       </a:blip>
       <a:srcRect/>
       <a:stretch>
         <a:fillRect/>
       </a:stretch>
     </p:blipFill>
     <p:spPr bwMode="auto">
       <a:xfrm>
         <a:off x="755576" y="692696"/>
         <a:ext cx="974725" cy="542925"/>
       </a:xfrm>
       <a:prstGeom prst="rect">
         <a:avLst/>
       </a:prstGeom>
       <a:noFill/>
       <a:extLst>
         <a:ext uri="{909E8E84-426E-40DD-AFC4-6F175D3DCCD1}">
           <a14:hiddenFill xmlns:a14="http://schemas.microsoft.com/office/drawing/2010/main">
             <a:solidFill>
               <a:srgbClr val="FFFFFF"/>
             </a:solidFill>
           </a14:hiddenFill>
         </a:ext>
       </a:extLst>
     </p:spPr>
   </p:pic>
*/

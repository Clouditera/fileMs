package picture

import "testing"

func TestPicture1(t *testing.T) {
	imageFile := "./static/uims.jpeg"
	ReadPic(imageFile)
}

func TestPicture(t *testing.T) {
	//WriteStringToFile("james.txt", ImageToChars(equalScaleImageFromWidth(*(GetImage("1.jpg")), 150, 150), "#$*@    "))
	//WriteStringToFile("./static/out/uims.txt", ImageToChars(equalScaleImageFromWidth(*(GetImage("./static/uims.jpeg")), 50, 50), "##  "))
	//WriteStringToFile("uims2.txt", ImageToChars(equalScaleImageFromWidth(*(GetImage("uims2.jpeg")), 50, 50), "##  "))
	WriteStringToFile("./static/out/logo.txt", ImageToChars(equalScaleImageFromWidth(*(GetImage("./static/logo.jpeg")), 100, 100), "##  "))
	//WriteStringToFile("smallSuperMan.txt", ImageToChars(equalScaleImageFromWidth(*(CrateImageFromString("1.png", "小超人\n BUG \n不会飞", 500, 500, 100)), 150, 150), "#$*@    "))
}

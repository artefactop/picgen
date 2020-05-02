# picgen
Picture generator service

> Note: Alpha stage not ready for production

![picgen example](https://picgen.xyz/256x158/ce3262/black.png?text=picgen&size=41)

http://localhost:3001/400x200/white/ce3262.png?text=picgen&size=65

It returns an image created with path values:
- Size: 400x200px 
- Background Color: white
- Label Color: #CE3262
- Image format: `image/png`
- Label Text: picgen
- Label Size: 65

You can set color with color keyword names as defined in [SVG 1.1](https://www.w3.org/TR/2003/REC-SVG11-20030114/types.html#ColorKeywords) or with its hexadecimal RGB value.

Image format supported:
- png (default)
- jpeg (no recommended)

Other API format

http://localhost:3001?x=400x200&b=000&f=ce3262?t=picgen

It returns an image created with path values:
- Size: 400x200px 
- Background Color: #FFFFFF
- Label Color: #CE3262
- Image format: `image/png`
- Label Text: picgen
- Label Size: default (65)


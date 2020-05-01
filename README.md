# picgen
Picture generator service

http://localhost:3001/400x200/white/ce3262?label=picgen&size=65

It returns an `image/png` with:
- Size: 400x200px 
- Background Color: white
- Label Color: #CE3262
- Label Text: picgen
- Label Size: 65

You can set color with color names as defined in SVG 1.1 or with it hexadecimal RGB value.

> Note: Alpha stage not ready for production
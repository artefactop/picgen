# Picture Generator Service

Pure [GO](https://golang.org) service that returns images based on URL path values that you can use: as placeholders while designing a web, make dinamic banners for your website, make default users avatars as in Gmail. See on [Github](https://github.com/artefactop/picgen)

![picgen example](https://picgen.xyz/256x158/ce3262/black.png?text=picgen&size=41)

## Examples

Placeholder: https://picgen.xyz/350x200/fcc5dc/001b46

![Placeholder image](https://picgen.xyz/350x200/fcc5dc/001b46)

Banner: https://picgen.xyz/550x128/7dffcf/1000c3?text=Banner

![Banner image](https://picgen.xyz/550x128/7dffcf/1000c3?text=Banner)

Avatar: https://picgen.xyz/64/deepskyblue/floralwhite?text=P&size=40

![Avatar image](https://picgen.xyz/64/deepskyblue/floralwhite?text=P&size=40)


## How to use

Request an image with this URL format:

`https://picgen.xyz/{Size}/{Background Color}/{Label Color}{Image Format}?text={Label Text}&size={Label Size}`

```
<img src="https://picgen.xyz/350/darkturquoise/gold">
<img src="https://picgen.xyz/350x200/3366CC/fff">
<img src="https://picgen.xyz/350x200/3366CC/fff?text=Hello">
```

> Color can be represented with its Hexadecimal value `#CE3262` or color keyword names as defined in [SVG 1.1](https://www.w3.org/TR/2003/REC-SVG11-20030114/types.html#ColorKeywords)

| Field  | Required | Default | Values |
|---|---|---|---|
| Size | yes  | none | Set width and height in pixels: `400x200`<br> Indicate only one value for squared images: `200` |
| Background Color   | yes  | none  | Color as hex or keyword name |
| Label Color  | yes | none | Color as hex or keyword name |
| Image Format | no  | png  | `.png` <br> `.jpeg` |
| Label Text   | no  | Size of the image  | URL encoded text  |
| Label Size   | no  | 65   | Integer |

### Query params API format

`https://picgen.xyz?x={Size}&b={Background Color}&f={Label Color}?t={Label Text}`

| Field  | Required | Default | Values |
|---|---|---|---|
| Size | no  | `100x100` | Set width and height in pixels: `400x200`<br> Indicate only one value for squared images: `200` |
| Background Color   | no  | `00ADD8`  | Color as hex or keyword name |
| Label Color  | no | `FFFFFF` | Color as hex or keyword name |
| Image Format | no  | png  | Not supported |
| Label Text   | no  | Size of the image  | URL encoded text  |
| Label Size   | no  | 65   | Not supported |

## Configure your service

You can fork this project and run it in your own GCP project. 
This project is ready to deploy on [Google Cloud Run](https://cloud.google.com/run). Configuration files are under [build](https://github.com/artefactop/picgen/tree/master/build) directory.

- Dockerfile: Service Dockerfile 
- ci/cloudbuild: 
    - Build service binary
    - Build docker image
    - Push docker image to Google Cloud Registry
    - Deploy docker image to Google Cloud Run

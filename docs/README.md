# Picture Generator Service

Pure [GO](https://golang.org) service that returns images that you can use: as placeholders while designing a web, make dinamic banners for your website, make default users avatars as in Gmail. See on [Github](https://github.com/artefactop/picgen)

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

> Color: Color can be represented with its Hexadecimal value `#CE3262` or color keyword names as defined in [SVG 1.1](https://www.w3.org/TR/2003/REC-SVG11-20030114/types.html#ColorKeywords)

- Size `required`: width x height in pixels examples: `400x200`, indicate only one value for squared images `200`
- Background Color `required`: Background color of image
- Label Color `required`: Label color
- Image Format: `.png` (default) or `.jpeg` (no recommended)
- Label Text: URL encoded text, by default it returns the size of image
- Label Size: Size of label, by default 65

### Query params API format

https://picgen.xyz?x={Size}&b={Background Color}&f={Label Color}?t={Label Text}

It returns an image created with path values:
- Size: `400x200`, defaults to `100x100`
- Background Color: `#FFFFFF`, defaults to `64C8C8`
- Label Color: `#CE3262`, defaults to `64C8C8`
- Image Format: Always `image/png`
- Label Text: URL encoded text, defaults to `#`
- Label Size: Defaults to 65

## Configure your service

You can fork this project and run it in your own GCP project. 
This project is ready to deploy on [Google Cloud Run](https://cloud.google.com/run). Configuration files are under [build](https://github.com/artefactop/picgen/tree/master/build) directory.

- Dockerfile: Service Dockerfile 
- ci/cloudbuild: 
    - Build service binary
    - Build docker image
    - Push docker image to Google Registry
    - Deploy docker image to Google Cloud Run

import Fastify from 'fastify'
import { chromium } from 'playwright'

const fastify = Fastify({ logger: true })

interface RenderRequest {
    templateJson: {
        pages: Array<{
            elements: Array<{
                id: string
                type: string
                x: number
                y: number
                width: number
                height: number
                text?: string
                src?: string
                style?: Record<string, any>
            }>
        }>
    }
}

// Helper to convert template JSON to HTML
// Helper to separate Layout vs Document logic
function htmlFromTipTap(node: any): string {
    if (!node) return ''

    if (node.type === 'doc') {
        return node.content.map((c: any) => htmlFromTipTap(c)).join('')
    }

    if (node.type === 'paragraph') {
        return `<p>${node.content ? node.content.map((c: any) => htmlFromTipTap(c)).join('') : '<br>'}</p>`
    }

    if (node.type === 'text') {
        let text = node.text
        if (node.marks) {
            node.marks.forEach((mark: any) => {
                if (mark.type === 'bold') text = `<b>${text}</b>`
                if (mark.type === 'italic') text = `<i>${text}</i>`
            })
        }
        return text
    }

    if (node.type === 'heading') {
        const level = node.attrs?.level || 1
        return `<h${level}>${node.content ? node.content.map((c: any) => htmlFromTipTap(c)).join('') : ''}</h${level}>`
    }

    if (node.type === 'bulletList') {
        return `<ul>${node.content ? node.content.map((c: any) => htmlFromTipTap(c)).join('') : ''}</ul>`
    }

    if (node.type === 'orderedList') {
        return `<ol>${node.content ? node.content.map((c: any) => htmlFromTipTap(c)).join('') : ''}</ol>`
    }

    if (node.type === 'listItem') {
        return `<li>${node.content ? node.content.map((c: any) => htmlFromTipTap(c)).join('') : ''}</li>`
    }

    if (node.type === 'variable') {
        return `<span style="background: #ede9fe; color: #5b21b6; padding: 2px 4px; border-radius: 4px; font-weight: 500;">{{ ${node.attrs?.label || 'var'} }}</span>`
    }

    return ''
}

function generateHTML(templateJson: any) {
    // Check if it's TipTap JSON (has type: 'doc')
    if (templateJson.type === 'doc') {
        const contentHtml = htmlFromTipTap(templateJson)
        return `
    <!DOCTYPE html>
    <html>
      <head>
        <style>
          @page { size: A4; margin: 2cm; }
          body { font-family: sans-serif; font-size: 12pt; line-height: 1.5; color: #333; }
          p { margin-bottom: 1em; }
          h1 { font-size: 24pt; margin-bottom: 0.5em; }
          h2 { font-size: 18pt; margin-bottom: 0.5em; }
          ul, ol { margin-bottom: 1em; padding-left: 2em; }
        </style>
      </head>
      <body>${contentHtml}</body>
    </html>
    `
    }

    // Default to Layout Editor (existing logic)
    const pageStyle = `
    @page { size: A4; margin: 0; }
    body { margin: 0; padding: 0; font-family: sans-serif; }
    .page { 
      width: 595pt; 
      height: 842pt; 
      position: relative; 
      page-break-after: always; 
      overflow: hidden;
      background: white;
    }
    .element { position: absolute; }
  `

    let pagesHtml = ''
    // Handle case where templateJson might be direct elements array (if migrated poorly) or object
    const pages = templateJson.pages || [{ elements: templateJson.elements || [] }]

    pages.forEach((page: any) => {
        let elementsHtml = ''
        if (page.elements) {
            page.elements.forEach((el: any) => {
                const style = `
             top: ${el.y}px;
             left: ${el.x}px;
             width: ${el.width}px;
             height: ${el.height}px;
             font-size: ${el.style?.fontSize || 16}px;
             color: ${el.style?.color || 'black'};
             white-space: pre-wrap;
           `

                if (el.type === 'text') {
                    elementsHtml += `<div class="element" style="${style}">${el.text || ''}</div>`
                } else if (el.type === 'image') {
                    elementsHtml += `<img class="element" src="${el.src}" style="${style}; object-fit: contain;" />`
                } else if (el.type === 'field') {
                    elementsHtml += `<div class="element" style="${style}; color: blue;">${el.text || '{{field}}'}</div>`
                }
            })
        }
        pagesHtml += `<div class="page">${elementsHtml}</div>`
    })

    return `
    <!DOCTYPE html>
    <html>
      <head><style>${pageStyle}</style></head>
      <body>${pagesHtml}</body>
    </html>
  `
}

fastify.post('/render', async (request, reply) => {
    const body = request.body as RenderRequest
    if (!body.templateJson) {
        return reply.code(400).send({ error: "templateJson required" })
    }

    const html = generateHTML(body.templateJson)

    // Launch Playwright
    const browser = await chromium.launch()
    const page = await browser.newPage()

    await page.setContent(html)

    // PDF Generation
    const pdfBuffer = await page.pdf({
        format: 'A4',
        printBackground: true
    })

    await browser.close()

    reply.header('Content-Type', 'application/pdf')
    reply.send(pdfBuffer)
})

const start = async () => {
    try {
        await fastify.listen({ port: 3001, host: '0.0.0.0' })
    } catch (err) {
        fastify.log.error(err)
        process.exit(1)
    }
}
start()

package endpoint

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mokiat/PipniAPI/internal/model/endpoint"
	"github.com/mokiat/PipniAPI/internal/widget"
	"github.com/mokiat/gog"
	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/ui"
	co "github.com/mokiat/lacking/ui/component"
	"github.com/mokiat/lacking/ui/layout"
	"github.com/mokiat/lacking/ui/mvc"
	"github.com/mokiat/lacking/ui/std"
)

var supportedMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

var Editor = mvc.EventListener(co.Define(&editorComponent{}))

type EditorData struct {
	EditorModel *endpoint.Editor
}

type editorComponent struct {
	co.BaseComponent

	mdlEditor *endpoint.Editor

	overlay co.Overlay
}

func (c *editorComponent) OnUpsert() {
	data := co.GetData[EditorData](c.Properties())
	c.mdlEditor = data.EditorModel
}

func (c *editorComponent) Render() co.Instance {
	methodItems := gog.Map(supportedMethods, func(method string) std.DropdownItem {
		return std.DropdownItem{
			Key:   method,
			Label: method,
		}
	})

	return co.New(std.Container, func() {
		co.WithLayoutData(c.Properties().LayoutData())
		co.WithData(std.ContainerData{
			BorderSize:  ui.UniformSpacing(1),
			BorderColor: opt.V(std.OutlineColor),
			Padding:     ui.UniformSpacing(2),
			Layout: layout.Frame(layout.FrameSettings{
				ContentSpacing: ui.Spacing{
					Top: 5,
				},
			}),
		})

		co.WithChild("uri-settings", co.New(std.Element, func() {
			co.WithLayoutData(layout.Data{
				VerticalAlignment: layout.VerticalAlignmentTop,
			})
			co.WithData(std.ElementData{
				Layout: layout.Frame(layout.FrameSettings{
					ContentSpacing: ui.SymmetricSpacing(5, 0),
				}),
			})

			co.WithChild("method", co.New(std.Dropdown, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentLeft,
					Width:               opt.V(100),
				})
				co.WithData(std.DropdownData{
					Items:       methodItems,
					SelectedKey: c.mdlEditor.Method(),
				})
				co.WithCallbackData(std.DropdownCallbackData{
					OnItemSelected: func(key any) {
						c.changeMethod(key.(string))
					},
				})
			}))

			co.WithChild("uri", co.New(widget.EditBox, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
				})
				co.WithData(widget.EditBoxData{
					Text: c.mdlEditor.URI(),
				})
				co.WithCallbackData(widget.EditBoxCallbackData{
					OnChange: func(text string) {
						c.changeURI(text)
					},
				})
			}))

			co.WithChild("go", co.New(std.Button, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentRight,
					Width:               opt.V(100),
				})
				co.WithData(std.ButtonData{
					Text: "GO",
				})
				co.WithCallbackData(std.ButtonCallbackData{
					OnClick: c.makeRequest,
				})
			}))
		}))

		co.WithChild("payload-settings", co.New(std.Container, func() {
			co.WithLayoutData(layout.Data{
				HorizontalAlignment: layout.HorizontalAlignmentCenter,
				VerticalAlignment:   layout.VerticalAlignmentCenter,
			})
			co.WithData(std.ContainerData{
				Layout: layout.Frame(),
			})

			co.WithChild("request", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentLeft,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
					Width:               opt.V(2000), // FIXME: Use different layout
				})
				co.WithData(std.ContainerData{
					Padding:     ui.UniformSpacing(1),
					BorderColor: opt.V(std.OutlineColor),
					BorderSize:  ui.UniformSpacing(1),
					Layout:      layout.Frame(),
				})

				co.WithChild("toolbar", co.New(std.Toolbar, func() {
					co.WithLayoutData(layout.Data{
						VerticalAlignment: layout.VerticalAlignmentTop,
					})
					co.WithData(std.ToolbarData{
						Positioning: std.ToolbarPositioningTop,
					})

					co.WithChild("body", co.New(std.ToolbarButton, func() {
						co.WithData(std.ToolbarButtonData{
							Icon:     co.OpenImage(c.Scope(), "images/data.png"),
							Text:     "Body",
							Selected: c.mdlEditor.RequestTab() == endpoint.EditorTabBody,
						})
						co.WithCallbackData(std.ToolbarButtonCallbackData{
							OnClick: func() {
								c.changeRequestTab(endpoint.EditorTabBody)
							},
						})
					}))

					co.WithChild("headers", co.New(std.ToolbarButton, func() {
						co.WithData(std.ToolbarButtonData{
							Icon:     co.OpenImage(c.Scope(), "images/headers.png"),
							Text:     "Headers",
							Selected: c.mdlEditor.RequestTab() == endpoint.EditorTabHeaders,
						})
						co.WithCallbackData(std.ToolbarButtonCallbackData{
							OnClick: func() {
								c.changeRequestTab(endpoint.EditorTabHeaders)
							},
						})
					}))
				}))

				switch c.mdlEditor.RequestTab() {
				case endpoint.EditorTabBody:
					co.WithChild("body", co.New(RequestBody, func() {
						co.WithData(RequestBodyData{
							Text: c.mdlEditor.RequestBody(),
						})
						co.WithCallbackData(RequestBodyCallbackData{
							OnChange: func(text string) {
								c.mdlEditor.SetRequestBody(text)
							},
						})
					}))

				case endpoint.EditorTabHeaders:
					co.WithChild("headers", co.New(RequestHeaders, func() {
						co.WithData(RequestHeadersData{
							EditorModel: c.mdlEditor,
						})
					}))
				}
			}))

			co.WithChild("response", co.New(std.Container, func() {
				co.WithLayoutData(layout.Data{
					HorizontalAlignment: layout.HorizontalAlignmentCenter,
					VerticalAlignment:   layout.VerticalAlignmentCenter,
				})
				co.WithData(std.ContainerData{
					Padding:     ui.UniformSpacing(1),
					BorderColor: opt.V(std.OutlineColor),
					BorderSize:  ui.UniformSpacing(1),
					Layout:      layout.Frame(),
				})

				co.WithChild("toolbar", co.New(std.Toolbar, func() {
					co.WithLayoutData(layout.Data{
						VerticalAlignment: layout.VerticalAlignmentTop,
					})
					co.WithData(std.ToolbarData{
						Positioning: std.ToolbarPositioningTop,
					})

					co.WithChild("body", co.New(std.ToolbarButton, func() {
						co.WithData(std.ToolbarButtonData{
							Icon:     co.OpenImage(c.Scope(), "images/data.png"),
							Text:     "Body",
							Selected: c.mdlEditor.ResponseTab() == endpoint.EditorTabBody,
						})
						co.WithCallbackData(std.ToolbarButtonCallbackData{
							OnClick: func() {
								c.changeResponseTab(endpoint.EditorTabBody)
							},
						})
					}))

					co.WithChild("headers", co.New(std.ToolbarButton, func() {
						co.WithData(std.ToolbarButtonData{
							Icon:     co.OpenImage(c.Scope(), "images/headers.png"),
							Text:     "Headers",
							Selected: c.mdlEditor.ResponseTab() == endpoint.EditorTabHeaders,
						})
						co.WithCallbackData(std.ToolbarButtonCallbackData{
							OnClick: func() {
								c.changeResponseTab(endpoint.EditorTabHeaders)
							},
						})
					}))

					co.WithChild("stats", co.New(std.ToolbarButton, func() {
						co.WithData(std.ToolbarButtonData{
							Icon:     co.OpenImage(c.Scope(), "images/stats.png"),
							Text:     "Stats",
							Selected: c.mdlEditor.ResponseTab() == endpoint.EditorTabStats,
							Enabled:  opt.V(false), // TODO: To be added
						})
						co.WithCallbackData(std.ToolbarButtonCallbackData{
							OnClick: func() {
								c.changeResponseTab(endpoint.EditorTabStats)
							},
						})
					}))
				}))

				switch c.mdlEditor.ResponseTab() {
				case endpoint.EditorTabBody:
					co.WithChild("body", co.New(ResponseBody, func() {
						co.WithData(ResponseBodyData{
							Text: c.mdlEditor.ResponseBody(),
						})
					}))

				case endpoint.EditorTabHeaders:
					co.WithChild("headers", co.New(ResponseHeaders, func() {
						co.WithData(ResponseHeadersData{
							EditorModel: c.mdlEditor,
						})
					}))

				case endpoint.EditorTabStats:
					// TODO: To be added
				}
			}))
		}))
	})
}

func (c *editorComponent) OnEvent(event mvc.Event) {
	switch event.(type) {
	case endpoint.MethodChangedEvent:
		c.Invalidate()
	case endpoint.URIChangedEvent:
		c.Invalidate()
	case endpoint.RequestTabChangedEvent:
		c.Invalidate()
	case endpoint.RequestBodyChangedEvent:
		c.Invalidate()
	case endpoint.ResponseTabChangedEvent:
		c.Invalidate()
	case endpoint.ResponseBodyChangedEvent:
		c.Invalidate()
	}
}

func (c *editorComponent) changeMethod(method string) {
	c.mdlEditor.SetMethod(method)
}

func (c *editorComponent) changeURI(uri string) {
	c.mdlEditor.SetURI(uri)
}

func (c *editorComponent) changeRequestTab(tab endpoint.EditorTab) {
	c.mdlEditor.SetRequestTab(tab)
}

func (c *editorComponent) changeResponseTab(tab endpoint.EditorTab) {
	c.mdlEditor.SetResponseTab(tab)
}

func (c *editorComponent) makeRequest() {
	ctx, ctxCancel := context.WithCancel(context.Background())

	c.overlay = co.OpenOverlay(c.Scope(), co.New(ActionModal, func() {
		co.WithCallbackData(ActionModalCallbackData{
			OnCancel: func() {
				ctxCancel()
				c.overlay.Close()
			},
		})
	}))

	op := constructCall(c.createRequest())
	go func() {
		response, err := op(ctx)
		co.Schedule(c.Scope(), func() {
			if err := ctx.Err(); err != nil {
				return // the request was cancelled
			}
			c.handleResponse(response, err)
			c.overlay.Close()
		})
	}()
}

func (c *editorComponent) createRequest() *APIRequest {
	return &APIRequest{
		Method:  c.mdlEditor.Method(),
		URI:     c.mdlEditor.URI(),
		Headers: c.mdlEditor.HTTPRequestHeaders(),
		Body:    c.mdlEditor.RequestBody(),
	}
}

func (c *editorComponent) handleResponse(response *APIResponse, err error) {
	if err == nil {
		c.mdlEditor.SetResponseBody(response.Body)
		c.mdlEditor.SetHTTPResponseHeaders(response.Headers)
	} else {
		log.Warn("API call error: %v", err)
		co.OpenOverlay(c.Scope(), co.New(widget.NotificationModal, func() {
			co.WithData(widget.NotificationModalData{
				Icon: co.OpenImage(c.Scope(), "images/error.png"),
				Text: fmt.Sprintf("Request error.\n\n%v", err.Error()),
			})
		}))
	}
}

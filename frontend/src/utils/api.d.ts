/**
 * This file was auto-generated by openapi-typescript.
 * Do not make direct changes to the file.
 */

export interface paths {
  '/account/approvals': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    get?: never
    put?: never
    /** Approve a auth request */
    post: operations['ApprovalsInterface_approve']
    delete?: never
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
  '/account/clients': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    /** Get all clients */
    get: operations['AccountClients_listClients']
    put?: never
    /** Create a new client */
    post: operations['AccountClients_createClient']
    delete?: never
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
  '/account/clients/{id}': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    /** Get a client */
    get: operations['AccountClients_getClient']
    put?: never
    /** Update a client */
    post: operations['AccountClients_updateClient']
    /** Delete a client */
    delete: operations['AccountClients_deleteClient']
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
  '/clients/{id}': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    /** Get a client */
    get: operations['ClientsInterface_getClient']
    put?: never
    post?: never
    delete?: never
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
  '/healthz': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    /** Check health */
    get: operations['Healthz_check']
    put?: never
    post?: never
    delete?: never
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
  '/oauth/authorize': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    /** Authorization Request */
    get: operations['OAuthInterface_authorize']
    put?: never
    /** Authorization Request */
    post: operations['OAuthInterface_postAuthorize']
    delete?: never
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
  '/oauth/jwks': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    /** Get JSON Web Key Set */
    get: operations['OAuthInterface_getJwks']
    put?: never
    post?: never
    delete?: never
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
  '/oauth/token': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    get?: never
    put?: never
    /** Get token */
    post: operations['OAuthInterface_getToken']
    delete?: never
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
  '/session': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    /** Get session */
    get: operations['SessionInterface_me']
    put?: never
    /** Login */
    post: operations['SessionInterface_login']
    /** Logout */
    delete: operations['SessionInterface_logout']
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
  '/users/signup': {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    get?: never
    put?: never
    /** Signup */
    post: operations['UsersInterface_signup']
    delete?: never
    options?: never
    head?: never
    patch?: never
    trace?: never
  }
}
export type webhooks = Record<string, never>
export interface components {
  schemas: {
    'AccountClients.CreateClientReq': {
      name: string
      redirect_urls: string[]
    }
    'AccountClients.CreatedClient': {
      id: string
      name: string
      secret: string
      redirect_urls: string[]
    }
    'AccountClients.UpdateClientReq': {
      name: string
      redirect_urls: string[]
    }
    /** @enum {string} */
    'Approvals.ApproveErr': 'invalid_client' | 'invalid_scope'
    'Approvals.ApproveReq': {
      client_id: string
      scope: string
    }
    Client: {
      id: string
      name: string
      redirect_urls: string[]
    }
    /** @enum {string} */
    'Clients.GetClientErr': 'client_not_found'
    'Clients.GetClientRes': {
      client: components['schemas']['PublicClient']
    }
    /** @enum {string} */
    'OAuth.AuthorizeErr':
      | 'invalid_request'
      | 'unauthorized_client'
      | 'access_denied'
      | 'unsupported_response_type'
      | 'invalid_scope'
      | 'server_error'
      | 'temporarily_unavailable'
    'OAuth.AuthorizeReqMultiPart': {
      response_type: components['schemas']['OAuth.AuthorizeResponseType']
      client_id: string
      redirect_uri: string
      scope: string
      state?: string
    }
    /** @enum {string} */
    'OAuth.AuthorizeResponseType': 'code'
    /** @enum {string} */
    'OAuth.TokenCacheControlHeader': 'no-store'
    /** @enum {string} */
    'OAuth.TokenErr': 'invalid_request'
    /** @enum {string} */
    'OAuth.TokenGrantType': 'authorization_code'
    /** @enum {string} */
    'OAuth.TokenPragmaHeader': 'no-store'
    'OAuth.TokenReq': {
      grant_type: components['schemas']['OAuth.TokenGrantType']
      code: string
      redirect_uri: string
      client_id: string
    }
    'OAuth.TokenRes': {
      access_token: string
      token_type: components['schemas']['OAuth.TokenTokenType']
      /** Format: int32 */
      expires_in: number
      /** Id token, if scope includes 'openid' */
      id_token?: string
    }
    /** @enum {string} */
    'OAuth.TokenTokenType': 'Bearer'
    PublicClient: {
      id: string
      name: string
    }
    /** @enum {string} */
    'Session.LoginErr': 'invalid_name_or_password'
    'Session.LoginReq': {
      name: string
      password: string
    }
    /** @enum {string} */
    'Session.LogoutErr': 'not_logged_in'
    'Session.MeRes': {
      user: components['schemas']['User'] | null
    }
    User: {
      /** Format: int64 */
      id: number
      name: string
    }
    'Users.ReqSignup': {
      name: string
      password: string
      password_confirmation: string
    }
    /** @enum {string} */
    'Users.SignupErr':
      | 'name_length_not_enough'
      | 'name_already_used'
      | 'password_length_not_enough'
      | 'password_confirmation_not_match'
  }
  responses: never
  parameters: never
  requestBodies: never
  headers: never
  pathItems: never
}
export type $defs = Record<string, never>
export interface operations {
  ApprovalsInterface_approve: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody: {
      content: {
        'application/json': components['schemas']['Approvals.ApproveReq']
      }
    }
    responses: {
      /** @description There is no content to send for this request, but the headers may be useful.  */
      204: {
        headers: {
          [name: string]: unknown
        }
        content?: never
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: components['schemas']['Approvals.ApproveErr']
            error_description: string
          }
        }
      }
    }
  }
  AccountClients_listClients: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody?: never
    responses: {
      /** @description The request has succeeded. */
      200: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            clients: components['schemas']['Client'][]
          }
        }
      }
    }
  }
  AccountClients_createClient: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody: {
      content: {
        'application/json': {
          client: components['schemas']['AccountClients.CreateClientReq']
        }
      }
    }
    responses: {
      /** @description The request has succeeded and a new resource has been created as a result. */
      201: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            client: components['schemas']['AccountClients.CreatedClient']
          }
        }
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: string
          }
        }
      }
    }
  }
  AccountClients_getClient: {
    parameters: {
      query?: never
      header?: never
      path: {
        id: string
      }
      cookie?: never
    }
    requestBody?: never
    responses: {
      /** @description The request has succeeded. */
      200: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            client: components['schemas']['Client']
          }
        }
      }
      /** @description The server cannot find the requested resource. */
      404: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: string
          }
        }
      }
    }
  }
  AccountClients_updateClient: {
    parameters: {
      query?: never
      header?: never
      path: {
        id: string
      }
      cookie?: never
    }
    requestBody: {
      content: {
        'application/json': {
          client: components['schemas']['AccountClients.UpdateClientReq']
        }
      }
    }
    responses: {
      /** @description There is no content to send for this request, but the headers may be useful.  */
      204: {
        headers: {
          [name: string]: unknown
        }
        content?: never
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: string
          }
        }
      }
    }
  }
  AccountClients_deleteClient: {
    parameters: {
      query?: never
      header?: never
      path: {
        id: string
      }
      cookie?: never
    }
    requestBody?: never
    responses: {
      /** @description There is no content to send for this request, but the headers may be useful.  */
      204: {
        headers: {
          [name: string]: unknown
        }
        content?: never
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: string
          }
        }
      }
    }
  }
  ClientsInterface_getClient: {
    parameters: {
      query?: never
      header?: never
      path: {
        id: string
      }
      cookie?: never
    }
    requestBody?: never
    responses: {
      /** @description The request has succeeded. */
      200: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': components['schemas']['Clients.GetClientRes']
        }
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: components['schemas']['Clients.GetClientErr']
            error_description: string
          }
        }
      }
    }
  }
  Healthz_check: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody?: never
    responses: {
      /** @description The request has succeeded. */
      200: {
        headers: {
          [name: string]: unknown
        }
        content?: never
      }
    }
  }
  OAuthInterface_authorize: {
    parameters: {
      query: {
        response_type: components['schemas']['OAuth.AuthorizeResponseType']
        client_id: string
        redirect_uri: string
        scope: string
        state?: string
      }
      header?: never
      path?: never
      cookie?: never
    }
    requestBody?: never
    responses: {
      /** @description There is no content to send for this request, but the headers may be useful.  */
      204: {
        headers: {
          [name: string]: unknown
        }
        content?: never
      }
      /** @description Redirection */
      302: {
        headers: {
          location: string
          [name: string]: unknown
        }
        content?: never
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: components['schemas']['OAuth.AuthorizeErr']
            error_description: string
            state?: string
          }
        }
      }
    }
  }
  OAuthInterface_postAuthorize: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody: {
      content: {
        'multipart/form-data': components['schemas']['OAuth.AuthorizeReqMultiPart']
      }
    }
    responses: {
      /** @description There is no content to send for this request, but the headers may be useful.  */
      204: {
        headers: {
          [name: string]: unknown
        }
        content?: never
      }
      /** @description Redirection */
      302: {
        headers: {
          location: string
          [name: string]: unknown
        }
        content?: never
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: components['schemas']['OAuth.AuthorizeErr']
            error_description: string
            state?: string
          }
        }
      }
    }
  }
  OAuthInterface_getJwks: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody?: never
    responses: {
      /** @description The request has succeeded. */
      200: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            [key: string]: unknown
          }
        }
      }
    }
  }
  OAuthInterface_getToken: {
    parameters: {
      query?: never
      header: {
        authorization: string
      }
      path?: never
      cookie?: never
    }
    requestBody: {
      content: {
        'application/x-www-form-urlencoded': components['schemas']['OAuth.TokenReq']
      }
    }
    responses: {
      /** @description The request has succeeded. */
      200: {
        headers: {
          'cache-control': components['schemas']['OAuth.TokenCacheControlHeader']
          pragma: components['schemas']['OAuth.TokenPragmaHeader']
          [name: string]: unknown
        }
        content: {
          'application/json': components['schemas']['OAuth.TokenRes']
        }
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: components['schemas']['OAuth.TokenErr']
            error_description: string
          }
        }
      }
    }
  }
  SessionInterface_me: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody?: never
    responses: {
      /** @description The request has succeeded. */
      200: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': components['schemas']['Session.MeRes']
        }
      }
    }
  }
  SessionInterface_login: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody: {
      content: {
        'application/json': components['schemas']['Session.LoginReq']
      }
    }
    responses: {
      /** @description There is no content to send for this request, but the headers may be useful.  */
      204: {
        headers: {
          'set-cookie': string
          [name: string]: unknown
        }
        content?: never
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: components['schemas']['Session.LoginErr']
            error_description: string
          }
        }
      }
    }
  }
  SessionInterface_logout: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody?: never
    responses: {
      /** @description There is no content to send for this request, but the headers may be useful.  */
      204: {
        headers: {
          'set-cookie': string
          [name: string]: unknown
        }
        content?: never
      }
    }
  }
  UsersInterface_signup: {
    parameters: {
      query?: never
      header?: never
      path?: never
      cookie?: never
    }
    requestBody: {
      content: {
        'application/json': components['schemas']['Users.ReqSignup']
      }
    }
    responses: {
      /** @description There is no content to send for this request, but the headers may be useful.  */
      204: {
        headers: {
          'set-cookie': string
          [name: string]: unknown
        }
        content?: never
      }
      /** @description The server could not understand the request due to invalid syntax. */
      400: {
        headers: {
          [name: string]: unknown
        }
        content: {
          'application/json': {
            error: components['schemas']['Users.SignupErr']
            error_description: string
          }
        }
      }
    }
  }
}

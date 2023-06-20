import { Entity, Host } from "./model"

export const fetchWithTimeout = async (
    url: RequestInfo,
    init: RequestInit,
    timeoutMs = 15 * 1000
): Promise<Response> => {
    const controller = new AbortController()
    const clientTimeout = setTimeout(() => {
        controller.abort()
    }, timeoutMs)

    try {
        const reqConfig: RequestInit = { ...init, signal: controller.signal }
        const res = await fetch(url, reqConfig)
        if (!res.ok) {
            const description = `status code:${res.status}`
            Promise.reject(new Error(description))
        }

        return res
    } catch (e: unknown) {
        if (e instanceof Error) {
            return await Promise.reject(new Error(`${e.name}: ${e.message}`))
        } else {
            return await Promise.reject(new Error('fetch failed with unknown error'))
        }
    } finally {
        clearTimeout(clientTimeout)
    }
}

export const getJWT = async (clientSignedToken: string): Promise<string> => {
    const requestOptions = {
        method: 'GET',
        headers: { authentication: clientSignedToken }
    }
    return await fetchWithTimeout(`/api/v1/auth/claim`, requestOptions)
        .then(async (res) => await res.json())
        .then((data) => {
            return data.jwt
        })
}

export const readEntity = async (ccid: string): Promise<Entity | undefined> => {
    return await  fetch(`api/v1/entity/${ccid}`, {
        method: 'GET',
        headers: {}
    }).then(async (res) => {
        const entity = await res.json()
        if (!entity || entity.ccaddr === '') {
            return undefined
        }
        return entity
    })
}

export const fetchWithCredential = async (jwt: string, url: RequestInfo, init: RequestInit, timeoutMs?: number): Promise<Response> => {
    const requestInit = {
        ...init,
        headers: {
            ...init.headers,
            authentication: 'Bearer ' + jwt
        }
    }
    return await fetchWithTimeout(url, requestInit, timeoutMs)
}


export const getHosts = async (): Promise<Host[]> => {
    return await fetch(`/api/v1/host/list`).then(async (data) => {
        return await data.json()
    })
}

export const getEntities = async (): Promise<Entity[]> => {
    return await fetch(`/api/v1/entity/list`).then(async (data) => {
        return await data.json()
    })
}

export const sayHello = async (jwt: string, remote: string): Promise<string> => {
    return await fetchWithCredential(jwt, `/api/v1/admin/sayhello/${remote}`, {}).then(async (data) => {
        return await data.json()
    })
}


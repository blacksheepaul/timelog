const bufferToBase64Url = (buffer: ArrayBuffer) => {
  const bytes = new Uint8Array(buffer)
  let binary = ''
  bytes.forEach(byte => {
    binary += String.fromCharCode(byte)
  })
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

const base64UrlToBuffer = (base64url: string) => {
  const padding = '='.repeat((4 - (base64url.length % 4)) % 4)
  const base64 = (base64url + padding).replace(/-/g, '+').replace(/_/g, '/')
  const binary = atob(base64)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i += 1) {
    bytes[i] = binary.charCodeAt(i)
  }
  return bytes.buffer
}

const mapCredentialCreationOptions = (options: any) => {
  const publicKey = options.publicKey || options
  return {
    ...publicKey,
    challenge: base64UrlToBuffer(publicKey.challenge),
    user: {
      ...publicKey.user,
      id:
        typeof publicKey.user.id === 'string'
          ? base64UrlToBuffer(publicKey.user.id)
          : publicKey.user.id,
    },
    excludeCredentials: (publicKey.excludeCredentials || []).map((cred: any) => ({
      ...cred,
      id: base64UrlToBuffer(cred.id),
    })),
  }
}

const mapCredentialRequestOptions = (options: any) => {
  const publicKey = options.publicKey || options
  return {
    ...publicKey,
    challenge: base64UrlToBuffer(publicKey.challenge),
    allowCredentials: (publicKey.allowCredentials || []).map((cred: any) => ({
      ...cred,
      id: base64UrlToBuffer(cred.id),
    })),
  }
}

const credentialToJSON = (credential: PublicKeyCredential) => {
  const response: any = credential.response
  const clientExtensionResults = credential.getClientExtensionResults()
  const common = {
    id: credential.id,
    rawId: bufferToBase64Url(credential.rawId),
    type: credential.type,
    clientExtensionResults,
    authenticatorAttachment: credential.authenticatorAttachment || undefined,
  }

  if (response instanceof AuthenticatorAttestationResponse) {
    return {
      ...common,
      response: {
        clientDataJSON: bufferToBase64Url(response.clientDataJSON),
        attestationObject: bufferToBase64Url(response.attestationObject),
        transports: response.getTransports ? response.getTransports() : undefined,
        authenticatorData: response.getAuthenticatorData
          ? bufferToBase64Url(response.getAuthenticatorData())
          : undefined,
        publicKey:
          response.getPublicKey && response.getPublicKey() !== null
            ? bufferToBase64Url(response.getPublicKey()!)
            : undefined,
        publicKeyAlgorithm: response.getPublicKeyAlgorithm
          ? response.getPublicKeyAlgorithm()
          : undefined,
      },
    }
  }

  const assertion = response as AuthenticatorAssertionResponse
  return {
    ...common,
    response: {
      clientDataJSON: bufferToBase64Url(assertion.clientDataJSON),
      authenticatorData: bufferToBase64Url(assertion.authenticatorData),
      signature: bufferToBase64Url(assertion.signature),
      userHandle: assertion.userHandle ? bufferToBase64Url(assertion.userHandle) : null,
    },
  }
}

export const isWebAuthnSupported = () =>
  typeof window !== 'undefined' && !!window.PublicKeyCredential && typeof navigator !== 'undefined'

export const beginRegistration = async (options: any) => {
  const publicKey = mapCredentialCreationOptions(options)
  const credential = (await navigator.credentials.create({ publicKey })) as PublicKeyCredential
  return credentialToJSON(credential)
}

export const beginLogin = async (options: any) => {
  const publicKey = mapCredentialRequestOptions(options)
  const credential = (await navigator.credentials.get({ publicKey })) as PublicKeyCredential
  return credentialToJSON(credential)
}

/**
 * TransportCodec defines the interface for encoding/decoding transport data.
 * This allows developers to plug in custom serialization strategies
 * (e.g., MessagePack, Protobuf, custom binary formats) instead of the default base64/JSON.
 */
export interface TransportCodec {
    /**
     * Decode response data from the transport format to a string
     * @param data - The encoded data (e.g., base64 string, binary array, etc.)
     * @param contentType - The content type of the response
     * @returns Decoded string data
     */
    decodeResponse(data: any, contentType: string): string;

    /**
     * Decode error data from the transport format to a string
     * @param data - The encoded error data
     * @returns Decoded error message
     */
    decodeError(data: any): string;

    /**
     * Encode request arguments for transport (optional)
     * @param args - The request arguments object
     * @returns Encoded data suitable for transport
     */
    encodeRequest?(args: any): any;
}

/**
 * Base64JSONCodec is the default codec that handles Go's JSON marshaling behavior.
 * Go JSON marshals []byte as base64 strings, so this codec decodes base64 back to UTF-8.
 */
export class Base64JSONCodec implements TransportCodec {
    /**
     * Decode base64-encoded response data to UTF-8 string
     */
    decodeResponse(data: any, contentType: string): string {
        if (!data) return '';

        try {
            // Decode base64 string to binary
            const binaryString = atob(data);
            const bytes = new Uint8Array(binaryString.length);
            for (let i = 0; i < binaryString.length; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }
            // Convert bytes to UTF-8 string
            return new TextDecoder().decode(bytes);
        } catch (err) {
            throw new Error(`Failed to decode response: ${err instanceof Error ? err.message : String(err)}`);
        }
    }

    /**
     * Decode base64-encoded error data to UTF-8 string
     */
    decodeError(data: any): string {
        if (!data) return 'Unknown error';

        try {
            return this.decodeResponse(data, 'text/plain');
        } catch {
            return 'Failed to decode error message';
        }
    }
}

/**
 * RawStringCodec handles responses that are already plain strings (no encoding).
 * Useful for transports that don't use base64 encoding.
 */
export class RawStringCodec implements TransportCodec {
    decodeResponse(data: any, contentType: string): string {
        return data ? String(data) : '';
    }

    decodeError(data: any): string {
        return data ? String(data) : 'Unknown error';
    }
}

/**
 * RawJSONCodec handles responses that are already JSON objects.
 * Useful for transports that send parsed JSON directly.
 */
export class RawJSONCodec implements TransportCodec {
    decodeResponse(data: any, contentType: string): string {
        if (!data) return '';
        return typeof data === 'string' ? data : JSON.stringify(data);
    }

    decodeError(data: any): string {
        if (!data) return 'Unknown error';
        return typeof data === 'string' ? data : JSON.stringify(data);
    }
}

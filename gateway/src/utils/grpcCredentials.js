import fs from 'fs';
import * as grpc from '@grpc/grpc-js';

const TLS_ENABLED = process.env.GRPC_TLS_ENABLED === 'true';
const TLS_CERT = process.env.GRPC_TLS_CERT || '/certs/server.crt';
const TLS_KEY = process.env.GRPC_TLS_KEY || '/certs/server.key';
const TLS_CA = process.env.GRPC_TLS_CA || '/certs/ca.crt';

/**
 * Returns gRPC channel credentials for client connections.
 * Uses TLS with mutual authentication when GRPC_TLS_ENABLED=true,
 * otherwise returns insecure credentials.
 */
export function getGrpcCredentials() {
  if (!TLS_ENABLED) {
    return grpc.credentials.createInsecure();
  }

  const ca = fs.readFileSync(TLS_CA);
  const cert = fs.readFileSync(TLS_CERT);
  const key = fs.readFileSync(TLS_KEY);

  return grpc.credentials.createSsl(ca, key, cert);
}

/**
 * Returns gRPC server credentials.
 * Uses TLS with client certificate verification when GRPC_TLS_ENABLED=true,
 * otherwise returns insecure server credentials.
 */
export function getServerCredentials() {
  if (!TLS_ENABLED) {
    return grpc.ServerCredentials.createInsecure();
  }

  const ca = fs.readFileSync(TLS_CA);
  const cert = fs.readFileSync(TLS_CERT);
  const key = fs.readFileSync(TLS_KEY);

  return grpc.ServerCredentials.createSsl(ca, [{ cert_chain: cert, private_key: key }], true);
}

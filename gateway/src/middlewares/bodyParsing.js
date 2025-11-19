import express from 'express';

const saveRawBody = (req, res, buf) => {
  if (buf && buf.length) {
    req.rawBody = buf.toString('utf8');
  }
};

/**
 * Body parsing middleware configuration
 */
export const bodyParsingMiddleware = (app) => {
  // Parse application/json
  app.use(express.json({ limit: '10mb', verify: saveRawBody }));

  // Parse application/x-www-form-urlencoded
  app.use(express.urlencoded({ extended: true, limit: '10mb', verify: saveRawBody }));

  // Parse text/plain
  app.use(express.text({ limit: '10mb', verify: saveRawBody }));

  // Parse application/octet-stream
  app.use(express.raw({ limit: '10mb', verify: saveRawBody }));
};

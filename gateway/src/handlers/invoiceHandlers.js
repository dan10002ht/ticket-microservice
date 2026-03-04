import grpcClients from '../grpc/clients.js';
import { sendSuccessResponse, createSimpleHandler } from '../utils/responseHandler.js';

// ============================================
// Invoice Operations
// ============================================

const getInvoice = async (req, res) => {
  const result = await grpcClients.invoiceService.GetInvoice({
    invoice_id: req.params.invoiceId,
    booking_id: '',
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const listInvoices = async (req, res) => {
  const result = await grpcClients.invoiceService.ListInvoices({
    user_id: req.user.id,
    status: req.query.status || '',
    page: parseInt(req.query.page) || 1,
    limit: parseInt(req.query.limit) || 20,
  });
  sendSuccessResponse(res, 200, result, req.correlationId);
};

const getInvoicePdf = async (req, res) => {
  const result = await grpcClients.invoiceService.GetInvoicePdf({
    invoice_id: req.params.invoiceId,
  });

  if (!result.success) {
    return sendSuccessResponse(res, 404, result, req.correlationId);
  }

  const pdfBuffer = Buffer.from(result.pdf_bytes);
  res.set({
    'Content-Type': result.content_type || 'application/pdf',
    'Content-Disposition': `attachment; filename="${result.filename || 'invoice.pdf'}"`,
    'Content-Length': pdfBuffer.length,
  });
  res.send(pdfBuffer);
};

// ============================================
// Exports
// ============================================

export const getInvoiceHandler = createSimpleHandler(getInvoice, 'invoice', 'getInvoice');
export const listInvoicesHandler = createSimpleHandler(listInvoices, 'invoice', 'listInvoices');
export const getInvoicePdfHandler = createSimpleHandler(getInvoicePdf, 'invoice', 'getInvoicePdf');

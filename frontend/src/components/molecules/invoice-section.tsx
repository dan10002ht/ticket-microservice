"use client";

import { FileText, Download, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { useInvoices, useDownloadInvoicePdf } from "@/lib/api/queries";
import { showToast } from "@/lib/toast";
import type { ApiError } from "@/lib/api/types/common";

interface InvoiceSectionProps {
  bookingId: string;
}

function formatPrice(price: number, currency = "VND") {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency,
    minimumFractionDigits: 0,
  }).format(price);
}

const statusVariant: Record<string, "default" | "secondary" | "destructive"> = {
  generated: "default",
  pending: "secondary",
  cancelled: "destructive",
};

export function InvoiceSection({ bookingId }: InvoiceSectionProps) {
  const { data } = useInvoices({ booking_id: bookingId });
  const downloadMutation = useDownloadInvoicePdf();

  const invoice = data?.items?.[0];

  if (!invoice) {
    return null;
  }

  const handleDownload = async () => {
    try {
      await downloadMutation.mutateAsync(invoice.id);
    } catch (err) {
      showToast.apiError(err as ApiError);
    }
  };

  return (
    <>
      <Separator />
      <div className="p-4 space-y-2">
        <div className="flex items-center justify-between">
          <h3 className="text-sm font-semibold flex items-center gap-2">
            <FileText className="h-4 w-4" />
            Invoice
          </h3>
          {invoice.status === "generated" && (
            <Button
              variant="outline"
              size="sm"
              onClick={handleDownload}
              disabled={downloadMutation.isPending}
            >
              {downloadMutation.isPending ? (
                <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />
              ) : (
                <Download className="mr-1.5 h-3.5 w-3.5" />
              )}
              Download PDF
            </Button>
          )}
        </div>
        <div className="grid gap-2 text-sm sm:grid-cols-2">
          <div>
            <span className="text-muted-foreground">Invoice #: </span>
            <span className="font-mono text-xs">{invoice.invoice_number}</span>
          </div>
          <div>
            <span className="text-muted-foreground">Status: </span>
            <Badge
              variant={statusVariant[invoice.status] ?? "secondary"}
              className="capitalize"
            >
              {invoice.status}
            </Badge>
          </div>
          <div>
            <span className="text-muted-foreground">Total: </span>
            <span className="font-medium">
              {formatPrice(invoice.total_amount, invoice.currency)}
            </span>
          </div>
          {invoice.issued_at && (
            <div>
              <span className="text-muted-foreground">Issued: </span>
              <span>
                {new Date(invoice.issued_at).toLocaleDateString("en-US")}
              </span>
            </div>
          )}
        </div>
      </div>
    </>
  );
}

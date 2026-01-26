'use client';

import { useState, useEffect } from 'react';
import { X, Copy, Check, AlertCircle } from 'lucide-react';

interface DomainVerificationModalProps {
  isOpen: boolean;
  onClose: () => void;
  domain: string;
  verificationToken: string;
  onVerify: () => Promise<void>;
}

export default function DomainVerificationModal({
  isOpen,
  onClose,
  domain,
  verificationToken,
  onVerify,
}: DomainVerificationModalProps) {
  const [copied, setCopied] = useState(false);
  const [verifying, setVerifying] = useState(false);
  const [verificationStatus, setVerificationStatus] = useState<'idle' | 'checking' | 'verified' | 'failed'>('idle');
  const [error, setError] = useState<string | null>(null);

  const txtRecord = `cortex-verification=${verificationToken}`;

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(txtRecord);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  const checkVerification = async () => {
    setVerifying(true);
    setError(null);
    setVerificationStatus('checking');
    
    try {
      await onVerify();
      setVerificationStatus('verified');
      // Close modal after 2 seconds on success
      setTimeout(() => {
        onClose();
        setVerificationStatus('idle');
      }, 2000);
    } catch (err: any) {
      setVerificationStatus('failed');
      setError(err.message || 'Verification failed. Please check your DNS records.');
    } finally {
      setVerifying(false);
    }
  };

  // Reset state when modal closes
  useEffect(() => {
    if (!isOpen) {
      setVerificationStatus('idle');
      setError(null);
    }
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl p-8 max-w-2xl w-full mx-4 shadow-2xl">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-white">Verify Domain Ownership</h2>
          <button
            onClick={onClose}
            className="text-slate-500 hover:text-white transition-colors p-2 hover:bg-slate-800 rounded-lg"
          >
            <X size={24} />
          </button>
        </div>

        {/* Instructions */}
        <div className="space-y-6">
          <div className="bg-slate-950 border border-slate-800 rounded-xl p-6">
            <h3 className="text-sm font-bold text-slate-400 uppercase tracking-widest mb-4">
              Step 1: Add DNS TXT Record
            </h3>
            <p className="text-slate-300 mb-4 text-sm leading-relaxed">
              Add the following DNS TXT record to your domain <span className="font-bold text-white">{domain}</span>:
            </p>
            
            {/* TXT Record Display */}
            <div className="bg-slate-900 border border-slate-700 rounded-lg p-4 mb-4">
              <div className="flex items-center justify-between gap-4">
                <code className="text-green-400 font-mono text-sm flex-1 break-all">
                  {txtRecord}
                </code>
                <button
                  onClick={copyToClipboard}
                  className="flex items-center gap-2 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors text-sm font-bold"
                >
                  {copied ? (
                    <>
                      <Check size={16} className="text-green-400" />
                      <span className="text-green-400">Copied!</span>
                    </>
                  ) : (
                    <>
                      <Copy size={16} />
                      <span>Copy</span>
                    </>
                  )}
                </button>
              </div>
            </div>

            <div className="bg-blue-500/10 border border-blue-500/20 rounded-lg p-4">
              <div className="flex gap-3">
                <AlertCircle size={20} className="text-blue-400 flex-shrink-0 mt-0.5" />
                <div className="text-sm text-blue-300">
                  <p className="font-bold mb-1">Important:</p>
                  <ul className="list-disc list-inside space-y-1 text-blue-200/80">
                    <li>DNS changes can take up to 48 hours to propagate</li>
                    <li>Most DNS providers update within 5-15 minutes</li>
                    <li>Make sure to add the record at the root domain level</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>

          {/* Step 2: Verify */}
          <div className="bg-slate-950 border border-slate-800 rounded-xl p-6">
            <h3 className="text-sm font-bold text-slate-400 uppercase tracking-widest mb-4">
              Step 2: Verify Domain
            </h3>
            <p className="text-slate-300 mb-4 text-sm">
              Once you've added the DNS record, click the button below to verify ownership.
            </p>

            {verificationStatus === 'verified' && (
              <div className="mb-4 bg-green-500/10 border border-green-500/20 rounded-lg p-4">
                <div className="flex items-center gap-3">
                  <Check size={20} className="text-green-400" />
                  <p className="text-green-400 font-bold">Domain verified successfully!</p>
                </div>
              </div>
            )}

            {verificationStatus === 'failed' && error && (
              <div className="mb-4 bg-red-500/10 border border-red-500/20 rounded-lg p-4">
                <div className="flex items-center gap-3">
                  <AlertCircle size={20} className="text-red-400" />
                  <p className="text-red-400 text-sm">{error}</p>
                </div>
              </div>
            )}

            <button
              onClick={checkVerification}
              disabled={verifying || verificationStatus === 'checking'}
              className="w-full bg-blue-600 hover:bg-blue-500 disabled:bg-slate-800 disabled:cursor-not-allowed text-white font-bold py-3 px-6 rounded-xl transition-all shadow-lg active:scale-95 flex items-center justify-center gap-2"
            >
              {verificationStatus === 'checking' ? (
                <>
                  <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                  <span>Checking DNS Records...</span>
                </>
              ) : verificationStatus === 'verified' ? (
                <>
                  <Check size={20} />
                  <span>Verified!</span>
                </>
              ) : (
                <span>Verify Domain</span>
              )}
            </button>
          </div>
        </div>

        {/* Footer */}
        <div className="mt-6 pt-6 border-t border-slate-800">
          <p className="text-xs text-slate-500 text-center">
            Need help? Check your DNS provider's documentation for adding TXT records.
          </p>
        </div>
      </div>
    </div>
  );
}

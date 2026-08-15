import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
import { cn, formatCurrency, formatDate } from './utils';

describe('Frontend Utils Unit Tests', () => {
  it('cn() correctly merges tailwind classes and filters falsey values', () => {
    const result = cn('px-2 py-1', false && 'hidden', true && 'text-white', 'px-4');
    assert.equal(result, 'py-1 text-white px-4');
  });

  it('formatCurrency() formats USD numbers properly', () => {
    const formattedZero = formatCurrency(0);
    assert.ok(formattedZero.includes('$0'));

    const formattedThousand = formatCurrency(1500);
    assert.ok(formattedThousand.includes('$1,500'));

    const formattedMillion = formatCurrency(2500000);
    assert.ok(formattedMillion.includes('$2,500,000'));
  });

  it('formatDate() formats valid ISO dates into Brazilian format or falls back', () => {
    const valid = formatDate('2026-08-15T14:30:00Z');
    assert.ok(valid.length > 0);

    const invalid = formatDate('not-a-valid-date');
    assert.equal(invalid, '');
  });
});

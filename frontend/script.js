class TaxCalculatorUI {
    constructor() {
        this.apiUrl = 'http://localhost:8080/api/v1';
        this.form = document.getElementById('tax-form');
        this.loadingEl = document.getElementById('loading');
        this.errorEl = document.getElementById('error');
        this.resultsEl = document.getElementById('results');
        
        this.bindEvents();
    }
    
    bindEvents() {
        this.form.addEventListener('submit', (e) => this.handleSubmit(e));
    }
    
    async handleSubmit(e) {
        e.preventDefault();
        
        const salary = parseFloat(document.getElementById('salary').value);
        const year = parseInt(document.getElementById('year').value);
        
        this.showLoading();
        this.hideError();
        this.hideResults();
        
        try {
            const result = await this.calculateTax(salary, year);
            this.showResults(result, salary);
        } catch (error) {
            this.showError(error.message);
        }
    }
    
    async calculateTax(salary, year) {
        const response = await fetch(`${this.apiUrl}/calculate-tax`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ salary, year })
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Calculation failed');
        }
        
        return await response.json();
    }
    
    showResults(result, salary) {
        this.hideLoading();
        
        const formatCurrency = (amount) => 
            new Intl.NumberFormat('en-CA', { 
                style: 'currency', 
                currency: 'CAD' 
            }).format(amount);
        
        const formatPercent = (rate) => 
            new Intl.NumberFormat('en-CA', { 
                style: 'percent', 
                minimumFractionDigits: 2,
                maximumFractionDigits: 2 
            }).format(rate);
        
        this.resultsEl.innerHTML = `
            <h2>Tax Calculation Results</h2>
            <div class="summary">
                <div class="summary-item">
                    <span class="label">Annual Salary:</span>
                    <span class="value">${formatCurrency(salary)}</span>
                </div>
                <div class="summary-item">
                    <span class="label">Total Tax Owed:</span>
                    <span class="value highlight">${formatCurrency(result.total_tax)}</span>
                </div>
                <div class="summary-item">
                    <span class="label">Effective Tax Rate:</span>
                    <span class="value">${formatPercent(result.effective_rate)}</span>
                </div>
                <div class="summary-item">
                    <span class="label">After-Tax Income:</span>
                    <span class="value">${formatCurrency(salary - result.total_tax)}</span>
                </div>
            </div>
            
            <h3>Tax Breakdown by Bracket</h3>
            <div class="breakdown">
                ${this.renderBracketBreakdown(result.bracket_breakdown)}
            </div>
        `;
        
        this.resultsEl.classList.remove('hidden');
    }
    
    renderBracketBreakdown(breakdown) {
        return breakdown.map(item => {
            const formatCurrency = (amount) => 
                new Intl.NumberFormat('en-CA', { 
                    style: 'currency', 
                    currency: 'CAD' 
                }).format(amount);
            
            const formatPercent = (rate) => 
                new Intl.NumberFormat('en-CA', { 
                    style: 'percent', 
                    minimumFractionDigits: 1,
                    maximumFractionDigits: 1 
                }).format(rate);
            
            const maxDisplay = item.bracket.max ? 
                formatCurrency(item.bracket.max) : 'and above';
            
            return `
                <div class="bracket-item">
                    <div class="bracket-range">
                        ${formatCurrency(item.bracket.min)} - ${maxDisplay}
                        <span class="bracket-rate">(${formatPercent(item.bracket.rate)})</span>
                    </div>
                    <div class="bracket-tax">${formatCurrency(item.tax_owed)}</div>
                </div>
            `;
        }).join('');
    }
    
    showLoading() {
        this.loadingEl.classList.remove('hidden');
        document.getElementById('submit-btn').disabled = true;
    }
    
    hideLoading() {
        this.loadingEl.classList.add('hidden');
        document.getElementById('submit-btn').disabled = false;
    }
    
    showError(message) {
        this.hideLoading();
        this.errorEl.textContent = message;
        this.errorEl.classList.remove('hidden');
    }
    
    hideError() {
        this.errorEl.classList.add('hidden');
    }
    
    hideResults() {
        this.resultsEl.classList.add('hidden');
    }
}

// Initialize the app
document.addEventListener('DOMContentLoaded', () => {
    new TaxCalculatorUI();
});
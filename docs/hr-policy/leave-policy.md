# Leave Policy (DRAFT — sample, to be tailored)

> **Status:** Dummy policy based on common practice in Indian IT companies (researched Sep 2026).
> It is a starting point, not legal advice. Leave laws differ by state (Shops & Establishments Act
> of the state where the company is registered) and by the new Labour Codes — have HR/legal confirm
> before adoption. Every number marked `[setting]` will be an admin-editable value in the app.

---

## 1. Purpose and scope

This policy explains the types of leave available to employees, how leave is earned, applied for,
approved, carried forward and encashed. It applies to all permanent employees. Interns and contract
staff follow their contract terms unless stated otherwise.

## 2. Leave year

- The leave year runs **1 April to 31 March** `[setting: leave year start month]`.
- Balances are reset / carried forward on the first day of the leave year.

## 3. Leave types at a glance

| Code | Leave type | Days per year | How it is credited | Carry forward | Encashable |
|------|-----------|---------------|--------------------|---------------|------------|
| CL | Casual Leave | 12 `[setting]` | Monthly: 1 day on the 1st of each month `[setting]` | No — lapses on 31 Mar | No |
| SL | Sick Leave | 8 `[setting]` | Full quota at start of year (pro-rata for joiners) | No — lapses on 31 Mar | No |
| EL | Earned / Privilege Leave | 18 `[setting]` | Monthly: 1.5 days on the 1st of each month | Yes, max 30 days total balance `[setting]` | Yes — at exit, and optionally above cap |
| ML | Maternity Leave | 182 days (26 weeks) | On request, per event | — | No |
| PL | Paternity Leave | 5 working days `[setting]` | On request, per event | — | No |
| BL | Bereavement Leave | 3 working days `[setting]` | On request, per event | — | No |
| MRL | Marriage Leave | 3 working days `[setting]` | Once during employment | — | No |
| CO | Compensatory Off | As earned | 1 day per full day worked on a weekly off / holiday, approved by manager | Must be used within 60 days `[setting]` | No |
| LOP | Loss of Pay | Max 15 days per year `[setting]` | Used only when no paid balance is available | — | — |

## 4. Rules for each leave type

### 4.1 Casual Leave (CL)
- For short, unplanned personal needs.
- Credited **1 day per month** `[setting]`, so an employee can only use what has been credited so far
  (e.g., by the end of June, at most 3 days).
- Maximum **2 consecutive days** per application `[setting]`. Longer absences must use EL.
- Apply at least **1 day in advance** where possible; same-day application allowed in emergencies.

### 4.2 Sick Leave (SL)
- For illness or medical appointments. May be applied on the day or after returning to work.
- A **medical certificate** is required for more than **2 consecutive days** `[setting]`.

### 4.3 Earned / Privilege Leave (EL)
- For planned leave such as vacations.
- Apply at least **7 days in advance** `[setting]`.
- Unused EL is carried forward, but the total balance cannot exceed **30 days** `[setting]`;
  days above the cap lapse (or are encashed, if the company chooses `[setting]`).
- Unused EL balance is **paid out at exit** (resignation, retirement, termination).

### 4.4 Maternity Leave (ML)
- **26 weeks paid** for the first two children; **12 weeks** from the third child onward
  (Maternity Benefit Act, as amended 2017).
- Eligibility: at least **80 days** worked in the 12 months before the expected delivery date.
- Up to 8 weeks may be taken before the expected delivery date.

### 4.5 Paternity Leave (PL)
- Not mandated by law for the private sector; this is a company benefit.
- **5 working days**, to be taken within **3 months** of the child's birth or adoption.

### 4.6 Bereavement Leave (BL)
- **3 working days** on the death of an immediate family member (spouse, child, parent, sibling,
  parent-in-law).

### 4.7 Marriage Leave (MRL)
- **3 working days**, once during employment, for the employee's own marriage.

### 4.8 Compensatory Off (CO)
- Earned when an employee works a full day on a weekly off or public holiday **with prior
  manager approval**.
- Must be used within **60 days** of being earned, otherwise it lapses.

### 4.9 Loss of Pay (LOP)
- Applies when the employee has no balance in the requested leave type, or for unauthorised absence.
- Salary is deducted for LOP days.
- LOP is **capped at 15 days per leave year** `[setting]`. The app will not accept an LOP request
  beyond the cap; the employee must contact HR, and only the HR admin can grant an exception.
- A single LOP request may not exceed **5 consecutive working days** `[setting]`.

## 5. How leave is counted

- **Working week:** Monday to Saturday; Sundays and the **2nd and 4th Saturday** of each month are weekly offs `[setting]`.
- **Half days** are allowed for CL, SL and EL `[setting]` (first half / second half).
- **Weekends and public holidays** falling inside a leave period are **not** counted as leave
  `[setting]`.
- **Sandwich rule:** not applied — a weekend between two leave days is not counted `[setting]`.
- **Negative balance:** not allowed. The app blocks a request that exceeds the available balance;
  the employee must apply the excess days separately as LOP (subject to the LOP cap).

## 5A. Settings available for every leave type

So the policy can change without code changes, each leave type has these admin-editable settings:

| Setting | Example (CL) |
|---------|--------------|
| Days per year (or "no quota") | 12 |
| Credit method: full quota at year start / monthly / on request per event | Monthly |
| Maximum consecutive days per application | 2 |
| Maximum days per year (hard cap, even if no quota) | 12 |
| Minimum advance notice (days) | 1 |
| Half day allowed | Yes |
| Attachment required after N consecutive days | — |
| Carry forward cap (0 = lapses) | 0 |
| Encashable at exit | No |
| Available during probation | Yes |
| Effective from date | 1 Apr 2027 |

## 6. New joiners, probation and exit

- **Pro-rata credit:** SL is credited in proportion to the months remaining in the leave year,
  including the month of joining (e.g., joining in October → 6/12 of the annual quota). CL and EL,
  being monthly, start from the next monthly credit.
- **During probation** (first 6 months `[setting]`): CL and SL may be used; EL accrues but can be
  used only after probation is confirmed `[setting]`.
- **At exit:** unused EL is encashed; CL, SL and CO lapse. Leave taken in excess of the pro-rata
  entitlement is recovered from final settlement.

## 7. Public holidays

- The company publishes a holiday list of **10 days** `[setting]` before the start of each leave year,
  including the national holidays (Republic Day, Independence Day, Gandhi Jayanti).
- The holiday list is maintained in the app and is used to calculate leave days.

## 8. Applying for and approving leave

1. The employee submits a leave request in the app: leave type, from/to dates, half-day option,
   reason, and attachment (if required, e.g., medical certificate).
2. The employee selects the **first approver** from the list of users with approval rights.
3. Each approver can **Approve & forward** to another approver, **Approve as final**, or **Reject**
   with a comment.
4. Leave balance is **deducted only on final approval**. While pending, the days are shown as
   "on hold" so they cannot be double-booked.
5. An employee can **cancel** a request that is still pending. Cancelling approved leave requires
   approval and restores the balance.
6. Every action is recorded in the request's history and cannot be edited.

## 9. Policy administration

- HR owns this policy and may revise it; changes apply from the stated effective date and do not
  alter balances or requests already approved.
- Values marked `[setting]` are configured by the HR admin in the application.

---

### Sources consulted
- Payoneer — Leave policy in India 2026: https://www.payoneer.com/en-in/resources/workforce-management/leave-policy/india/
- AYP Group — India Leave Policy Employer Guide 2026: https://ayp-group.com/blog/leave-policy-in-india
- WageIndicator — Earned & Casual Leave in India: https://wageindicator.org/en-in/work-in-india/labour-law/annual-leave-and-holidays/earned-casual-leave-in-india/
- SalaryBox — EL/CL/SL rules, accumulation & encashment 2026: https://salarybox.in/blog/earned-leave-casual-leave-sick-leave-in-india-rules-accumulation-encashment-explained-2026/
- greytHR — Maternity Leave in India 2026: https://www.greythr.com/blog/maternity-leave/
- Wisemonk — Paternity leave in India 2026: https://www.wisemonk.io/blogs/paternity-leave-in-india

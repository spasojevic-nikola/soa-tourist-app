import { TestBed } from '@angular/core/testing';

import { StakeholdersService } from './stakeholders.service';

describe('StakeholdersService', () => {
  let service: StakeholdersService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(StakeholdersService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});

import { Injectable } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Observable } from 'rxjs';
import { CreateKeyPointPayload } from '../../tour-keypoints/model/keypoint.model';
import { TourKeypointsComponent } from '../../tour-keypoints/tour-keypoints/tour-keypoints.component';

@Injectable({
  providedIn: 'root'
})
export class KeypointDialogService {

  constructor(private dialog: MatDialog) { }

  openKeypointDialog(latitude: number, longitude: number, order: number, existingKeyPoint?: CreateKeyPointPayload): Observable<CreateKeyPointPayload | undefined> {
    const dialogRef = this.dialog.open(TourKeypointsComponent, {
      width: '800px',
      height: 'auto',
      maxHeight: '90vh',
      data: { latitude, longitude, order, existingKeyPoint: existingKeyPoint },
      panelClass: 'keypoint-dialog-container'
    });

    return dialogRef.afterClosed();
  }
}